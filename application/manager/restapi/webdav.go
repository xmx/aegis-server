package restapi

import (
	"net/http"
	"path"

	"github.com/labstack/echo/v5"
	"golang.org/x/net/webdav"
)

func NewWebDAV(basepath, dir string) *WebDAV {
	prefix := path.Join(basepath, "/webdav")
	return &WebDAV{
		webDAV: &webdav.Handler{
			Prefix:     prefix,
			FileSystem: webdav.Dir(dir),
			LockSystem: webdav.NewMemLS(),
		},
	}
}

type WebDAV struct {
	webDAV *webdav.Handler
}

func (d *WebDAV) RegisterRoute(r *echo.Group) error {
	methods := []string{
		http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete,
		"LOCK", "UNLOCK", "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE",
	}

	r.Match(methods, "/webdav", d.browse)
	r.Match(methods, "/webdav/*", d.browse)

	return nil
}

func (d *WebDAV) browse(c *echo.Context) error {
	w, r := c.Response(), c.Request()
	d.webDAV.ServeHTTP(w, r)
	return nil
}
