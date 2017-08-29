package main

import (
	"github.com/antonkuzmenko/gogarin/pkg/satellite"
)

func FileCreated() {}
func FileUpdated() {}
func FileDeleted() {}

func main() {
	spaceCenter := satellite.NewSpaceCenter(satellite.SpaceCenterConfig{})
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
	sat.Start(spaceCenter)
	sat.Stop()
}
