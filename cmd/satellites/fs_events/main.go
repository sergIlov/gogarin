package main

import (
	"fmt"

	"github.com/antonkuzmenko/gogarin/pkg/satellite"
)

type RedisConnector struct {
}

func (c *RedisConnector) Register(s *satellite.Satellite) error {
	fmt.Printf("Satellite:\n%+v\n", s.Info)

	if len(s.Triggers) > 0 {
		fmt.Println("\nTriggers:")
	}
	for _, t := range s.Triggers {
		fmt.Printf("%+v\n", t.Info)
	}

	return nil
}

func FileCreated() {}

type FileCreatedConfig struct {
	Path      []string `json:"path" desc:"Path to the file or directory."`
	Recursive bool     `json:"recursive" desc:"Triggers when a new file is created n-tiers down the directory tree."`
}

func main() {
	sat := satellite.New(
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

	sat.Start(&RedisConnector{})
	sat.Stop()
}
