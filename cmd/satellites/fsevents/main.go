package main

import (
	"log"

	"github.com/antonkuzmenko/gogarin/pkg/satellite"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func FileCreated() {}

type FileCreatedConfig struct {
	Path      []string `json:"path" desc:"Path to the file or directory."`
	Recursive bool     `json:"recursive" desc:"Triggers when a new file is created n-tiers down the directory tree."`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var c satellite.Config
	err = envconfig.Process("gogarin_satellite", &c)
	if err != nil {
		log.Fatal(err)
	}

	rpc, err := satellite.NewRPC(c.RPC)
	if err != nil {
		log.Fatal(err)
	}

	sat := satellite.New(
		rpc,
		satellite.Info{
			Name:        "File System Events",
			Version:     "0.1.0-alpha",
			Description: "Provides a mechanism for monitoring file system events.",
		},
	)

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

	err = sat.Start()
	if err != nil {
		panic(err)
	}
	sat.Stop()
}
