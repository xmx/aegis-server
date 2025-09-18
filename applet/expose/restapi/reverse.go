package restapi

import (
	"net/http"
	"net/http/httputil"
	"path"
	"strings"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-control/library/httpnet"
)

func NewReverse(trip http.RoundTripper) *Reverse {
	prx := httpnet.NewReverse(trip)
	return &Reverse{
		prx: prx,
	}
}

type Reverse struct {
	prx *httputil.ReverseProxy
}

func (rvs *Reverse) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/reverse/agent/:id").Any(rvs.agent)
	r.Route("/reverse/agent/:id/*path").Any(rvs.agent)
	r.Route("/reverse/broker/:id").Any(rvs.broker)
	r.Route("/reverse/broker/:id/*path").Any(rvs.broker)
	return nil
}

func (rvs *Reverse) agent(c *ship.Context) error {
	id, pth := c.Param("id"), "/"+c.Param("path")
	//oid, _ := bson.ObjectIDFromHex(id)
	////

	w, r := c.Response(), c.Request()
	rawPath := r.URL.Path
	if pth != "/" && strings.HasSuffix(rawPath, "/") {
		pth += "/"
	}
	pth = path.Join("/api/reverse/agent/", id, pth)
	reqURL := transport.NewBrokerIDURL(id, pth)
	r.URL = reqURL
	r.Host = reqURL.Host

	rvs.prx.ServeHTTP(w, r)

	return nil
}

func (rvs *Reverse) broker(c *ship.Context) error {
	id, path := c.Param("id"), "/"+c.Param("path")
	w, r := c.Response(), c.Request()
	rawPath := r.URL.Path
	if path != "/" && strings.HasSuffix(rawPath, "/") {
		path += "/"
	}

	reqURL := transport.NewBrokerIDURL(id, path)
	r.URL = reqURL
	r.Host = reqURL.Host

	rvs.prx.ServeHTTP(w, r)

	return nil
}
