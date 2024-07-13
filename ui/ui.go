package ui

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var ui embed.FS

func UI() fs.FS {
	fsys, _ := fs.Sub(ui, "dist")

	return fsys
}
