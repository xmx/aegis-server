package restapi

import (
	"net/http"
	"strings"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewStatic(dir string) shipx.Controller {
	return &staticAPI{
		fs: http.FileServer(http.Dir(dir)),
	}
}

type staticAPI struct {
	fs http.Handler
}

func (api *staticAPI) Register(rt shipx.Router) error {
	anon := rt.Anon()
	anon.Route("/webui").
		GET(api.FS).
		HEAD(api.FS)
	anon.Route("/webui/*path").
		GET(api.FS).
		HEAD(api.FS)

	return nil
}

func (api *staticAPI) FS(c *ship.Context) error {
	w, r := c.Response(), c.Request()
	rawPath, argPath := r.URL.Path, c.Param("path")
	suffix := strings.HasSuffix(rawPath, "/")
	if argPath == "" && !suffix {
		return c.Redirect(http.StatusTemporaryRedirect, rawPath+"/")
	}
	if suffix {
		argPath += "/"
	}

	r.URL.Path = "/" + argPath
	api.fs.ServeHTTP(w, r)

	return nil
}
