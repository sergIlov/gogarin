package file_system

import "github.com/antonkuzmenko/gogarin/pkg/satellite"

const version = "0.1.0-alpha"

type FileSystem struct {
	satellite.Base
}

func (fs *FileSystem) Info() satellite.Info {
	return satellite.Info{
		Name:        "File System",
		Version:     version,
		Description: "File system events + CRUD for files, directories, and links.",
	}
}

func (fs *FileSystem) Triggers() []satellite.Trigger {
	return []satellite.Trigger{
		&FileCreated{},
	}
}
