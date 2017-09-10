package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonkuzmenko/gogarin/pkg/satellite"
	"github.com/antonkuzmenko/gogarin/pkg/transport"
	"github.com/antonkuzmenko/gogarin/pkg/transport/redis"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// set by release script, see Makefile
var (
	version   = "dev"
	commit    string
	buildTime string
)

type Config struct {
	RPC struct {
		Adapter             string `required:"true"`
		Redis               redis.Config
		PollTimeoutInMs     int `default:"2000"`
		ShutdownTimeoutInMs int `default:"30000"`
	}
	Logger string `default:"json"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var config Config
	err = envconfig.Process("gogarin_space_center", &config)
	if err != nil {
		panic(err)
	}

	logger := newLogger(config)
	conn := newConn(config, logger)

	enc := func(ctx context.Context, res interface{}) (response interface{}, err error) {
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(res)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	registerEndpoint := func(ctx context.Context, req interface{}) (res interface{}, err error) {
		i := req.(satellite.Info)
		level.Info(logger).Log("name", i.Name, "version", i.Version)
		return i, nil
	}

	registerDec := func(ctx context.Context, req interface{}) (res interface{}, err error) {
		r := req.([]byte)
		var i satellite.Info
		err = json.Unmarshal(r, &i)
		if err != nil {
			return nil, err
		}

		return i, nil
	}

	register := redis.NewServer(
		registerEndpoint, registerDec, enc,
		redis.ServerLogger(log.With(logger, "component", "redis.Server")),
	)
	server := transport.NewServer(
		conn,
		time.Duration(config.RPC.PollTimeoutInMs)*time.Millisecond,
		log.With(logger, "component", "transport.Server"),
	)
	server.Handle("satellite.register", register)
	go func() {
		er := server.Serve()
		if er != transport.ErrServerClosed {
			level.Error(log.With(logger, "component", "transport.Server")).Log("err", er, "context", "Serve")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	level.Info(logger).Log("sig", sig)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.RPC.ShutdownTimeoutInMs)*time.Millisecond,
	)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
}

const (
	redisRPC = "redis"
)

func newConn(c Config, l log.Logger) transport.Connection {
	if c.RPC.Adapter == redisRPC {
		return redis.New(c.RPC.Redis)
	}

	level.Error(l).Log("err", "invalid RPC.Adapter", "adapter", c.RPC.Adapter)
	os.Exit(1)
	return nil
}

const (
	jsonLogger   = "json"
	logfmtLogger = "logfmt"
)

func newLogger(c Config) log.Logger {
	var logger log.Logger

	switch c.Logger {
	case jsonLogger:
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	case logfmtLogger:
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	default:
		panic("invalid logger format: " + c.Logger)
	}

	logger = log.With(
		logger,
		"ts", log.DefaultTimestampUTC,
		"version", version,
		"commit", commit,
		"build_ts", buildTime,
	)
	return logger
}
