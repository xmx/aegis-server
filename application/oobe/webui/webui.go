package webui

import (
	"embed"
	"io/fs"
)

//go:embed dist
var embedFS embed.FS

func Dist() fs.FS {
	sub, _ := fs.Sub(embedFS, "dist")
	if sub == nil {
		sub = embedFS
	}

	return sub
}
