package webui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var embedFS embed.FS

func Dist() http.FileSystem {
	sub, _ := fs.Sub(embedFS, "dist")
	if sub == nil {
		sub = embedFS
	}

	return http.FS(sub)
}
