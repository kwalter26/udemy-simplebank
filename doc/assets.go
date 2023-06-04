package doc

import (
	"embed"
	"io/fs"
)

//go:embed all:swagger

var assets embed.FS

func Assets() (fs.FS, error) {
	return fs.Sub(assets, "swagger")
}
