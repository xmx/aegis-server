package restapi

import (
	"net/http"
	"strings"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/protocol/webfs"
)

func NewDAV(prefix, dir string) shipx.Router {
	prefix = strings.TrimRight(prefix, "/")
	const path = "/dav"
	prefix += path
	dav := webfs.DAV(prefix, dir)

	return &davAPI{
		path:   path,
		prefix: prefix,
		dav:    dav,
	}
}

type davAPI struct {
	path   string
	prefix string // HTTP path 前缀
	dav    http.Handler
}

func (api *davAPI) Route(r *ship.RouteGroupBuilder) error {
	methods := []string{
		http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete,
		"LOCK", "UNLOCK", "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE",
	}

	r.Route(api.path).Method(api.FS, methods...)
	r.Route(api.path+"/*wildcard").Method(api.FS, methods...)
	return nil
}

func (api *davAPI) FS(c *ship.Context) error {
	wildcard := c.Param("wildcard")
	w, r := c.ResponseWriter(), c.Request()
	path := r.URL.Path
	if wildcard == "" && !strings.HasSuffix(path, "/") {
		return c.Redirect(http.StatusPermanentRedirect, path+"/")
	}

	api.dav.ServeHTTP(w, r)

	return nil
}
