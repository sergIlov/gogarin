package main

import (
	"fmt"
	"os"

	"github.com/antonkuzmenko/gogarin/pkg/satellite"
	"github.com/antonkuzmenko/gogarin/pkg/schema"
	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func FileCreated() {}

type FileCreatedConfig struct {
	Path      []string `json:"path" desc:"Path to the file or directory."`
	Recursive bool     `json:"recursive" desc:"Triggers when a new file is created n-tiers down the directory tree."`
}

type AppendFileConfig struct {
	Path string `json:"path" desc:"Path to the file"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var c satellite.Config
	err = envconfig.Process("gogarin_satellite", &c)
	if err != nil {
		panic(err)
	}

	logger := newLogger(c)
	conn := satellite.NewConnection(c, logger)

	sat := satellite.New(
		conn,
		satellite.Info{
			Name:        "File System Events",
			Version:     "0.1.0-alpha",
			Description: "Provides a mechanism for monitoring file system events.",
		},
	)

	f := FileCreatedFields()
	fmt.Println(f)

	sat.AddTrigger(
		satellite.Trigger{
			Call: FileCreated,
			Info: satellite.AbilityInfo{
				Name:        "File Created",
				Description: "Triggers when a new file or directory is created.",
			},
			Config: FileCreatedConfig{},
			Validator: func(config interface{}) {

			},
		},
	)

	sat.AddAction(
		satellite.Action{
			Call: func() {},
			Info: satellite.AbilityInfo{
				Name:        "Append file",
				Description: "Append the specified file",
			},
			Config: AppendFileConfig{},
			Validator: func(config interface{}) {

			},
		},
	)

	err = sat.Start(c)
	if err != nil {
		panic(err)
	}
	sat.Stop()
}

const (
	jsonLogger   = "json"
	logfmtLogger = "logfmt"
)

func newLogger(c satellite.Config) log.Logger {
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

func FileCreatedFields() schema.Fields {
	return schema.Fields{
		"file": &schema.Field{
			Name:        "File",
			Type:        schema.Object,
			Description: "Created file or directory",
			Fields: schema.Fields{
				"name": &schema.Field{
					Name:        "File.Name",
					Type:        schema.String,
					Description: "Name of created file or directory",
				},
				"type": &schema.Field{
					Name:        "File.Type",
					Type:        schema.String,
					Description: "Name of created object (file or directory)",
				},
			},
		},
	}
}
