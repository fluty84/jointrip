package main

import (
	"embed"
	"io/fs"
)

// Embed the React build files
//go:embed web/dist/*
var webFiles embed.FS

// GetWebFS returns the embedded web filesystem
func GetWebFS() fs.FS {
	// Strip the "web/dist" prefix to serve files from root
	webFS, err := fs.Sub(webFiles, "web/dist")
	if err != nil {
		panic(err)
	}
	return webFS
}
