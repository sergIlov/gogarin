package satellite

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
	"github.com/antonkuzmenko/gogarin/pkg/transport/redis"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// RPCConfig contains configurations for RPC transport.
type RPCConfig struct {
	Adapter              string `required:"true"`
	Redis                redis.Config
	RegisterTimeoutInSec int `default:"10000"`
}

const (
	redisRPC = "redis"
)

// NewConnection creates new transport.Connection.
func NewConnection(c Config, logger log.Logger) transport.Connection {
	if c.RPC.Adapter == redisRPC {
		return redis.New(c.RPC.Redis)
	}

	level.Error(logger).Log("err", "invalid RPC.Adapter", "adapter", c.RPC.Adapter)
	os.Exit(1)
	return nil
}

func createRegisterEndpoint(c Config, conn transport.Connection) endpoint.Endpoint {
	return transport.NewClient(
		conn,
		"satellite.register",
		time.Duration(c.RPC.RegisterTimeoutInSec)*time.Second,
		func(ctx context.Context, data interface{}) (request interface{}, err error) {
			var buf bytes.Buffer
			err = json.NewEncoder(&buf).Encode(data)
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		},
		func(ctx context.Context, data interface{}) (response interface{}, err error) {
			r := data.([]byte)
			var i Info
			err = json.Unmarshal(r, &i)
			if err != nil {
				return nil, err
			}

			return i, nil
		},
	).Endpoint()
}
