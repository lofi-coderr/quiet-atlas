// Package web provides the compiled browser application.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var files embed.FS

// Assets removes the build directory prefix from the embedded filesystem.
func Assets() (fs.FS, error) { return fs.Sub(files, "dist") }
