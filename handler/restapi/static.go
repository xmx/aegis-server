package restapi

import (
	"net/http"
	"path"
	"strings"

	"github.com/xgfone/ship/v5"
)

func NewStatic(secPath, dir string) *Static {
	return &Static{
		path: secPath,
		fs:   http.FileServer(http.Dir(dir)),
	}
}

type Static struct {
	path string
	fs   http.Handler
}

func (api *Static) Route(r *ship.RouteGroupBuilder) error {
	r.Route(api.path).GET(api.serve).HEAD(api.serve)
	r.Route(path.Join(api.path, "*path")).GET(api.serve).HEAD(api.serve)

	return nil
}

func (api *Static) serve(c *ship.Context) error {
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
