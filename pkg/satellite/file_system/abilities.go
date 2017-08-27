package file_system

import (
	"github.com/antonkuzmenko/gogarin/pkg/satellite"
	_ "github.com/fsnotify/fsnotify"
)

type FileCreated struct {
	satellite.Trigger
}

func (fc *FileCreated) Info() satellite.AbilityInfo {
	return satellite.AbilityInfo{
		Name:        "File Created",
		Description: "Triggers when a new file is created.",
		Push:        true,
	}
}

func (fc *FileCreated) Messages() ([]*satellite.Message, error) {
	//w, err := fsnotify.NewWatcher()
	var ms []*satellite.Message
	return ms, nil
}
