package restapi

import (
	"net/http"
	"path"
	"strings"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewStatic(secPath, dir string) shipx.Router {
	return &staticAPI{
		path: secPath,
		fs:   http.FileServer(http.Dir(dir)),
	}
}

type staticAPI struct {
	path string
	fs   http.Handler
}

func (api *staticAPI) Route(r *ship.RouteGroupBuilder) error {
	r.Route(api.path).GET(api.FS).HEAD(api.FS)
	r.Route(path.Join(api.path, "*path")).GET(api.FS).HEAD(api.FS)

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
