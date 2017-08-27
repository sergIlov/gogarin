package main

import (
	"github.com/antonkuzmenko/gogarin/pkg/satellite"
	"github.com/antonkuzmenko/gogarin/pkg/satellite/file_system"
)

func main() {
	fs := &file_system.FileSystem{}
	satellite.Register(fs)
}
