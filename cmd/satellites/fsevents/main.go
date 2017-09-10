package main

import (
	"log"
	"github.com/antonkuzmenko/gogarin/pkg/satellite"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"fmt"
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

	f := FileCreatedFields()
	fmt.Print(f)
	
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

func FileCreatedFields() satellite.AbilityFields {
	return satellite.AbilityFields{
		"file": &satellite.AbilityField{
			Name:        "File",
			Type:        satellite.Object,
			Description: "Created file or directory",
			Fields: satellite.AbilityFields{
				"name": &satellite.AbilityField{
					Name:        "File.Name",
					Type:        satellite.String,
					Description: "Name of created file or directory",
				},
				"type": &satellite.AbilityField{
					Name:        "File.Type",
					Type:        satellite.String,
					Description: "Name of created object (file or directory)",
				},
			},
		},
	}
}
