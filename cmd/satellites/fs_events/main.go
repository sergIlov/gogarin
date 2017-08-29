package main

import (
	"fmt"

	"github.com/antonkuzmenko/gogarin/pkg/satellite"
)

type RabbitMQConnector struct {
}

func (c *RabbitMQConnector) Register(s *satellite.Satellite) error {
	fmt.Printf("Satellite:\n%+v\n", s.Info)

	if len(s.Triggers) > 0 {
		fmt.Println("\nTriggers:")
	}
	for i := range s.Triggers {
		fmt.Printf("%+v\n", i)
	}

	return nil
}

func FileCreated() {}
func FileUpdated() {}
func FileDeleted() {}

func main() {
	sat := satellite.New(
		satellite.Info{
			Name:        "File System Events",
			Version:     "0.1.0-alpha",
			Description: "Provides a mechanism for monitoring file system events.",
		},
	)

	sat.AddTrigger(
		FileCreated,
		satellite.AbilityInfo{
			Name:        "File Created",
			Description: "Triggers when a new file or directory is created.",
		},
	)
	sat.AddTrigger(
		FileUpdated,
		satellite.AbilityInfo{
			Name:        "File Updated",
			Description: "Triggers when a file or directory is changed.",
		},
	)
	sat.AddTrigger(
		FileDeleted,
		satellite.AbilityInfo{
			Name:        "File Deleted",
			Description: "Triggers when a file or directory is deleted.",
		},
	)

	connector := &RabbitMQConnector{}
	sat.Start(connector)
	sat.Stop()
}
