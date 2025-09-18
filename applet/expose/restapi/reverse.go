package restapi

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/library/httpnet"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewReverse(trip http.RoundTripper, repo repository.All) *Reverse {
	prx := httpnet.NewReverse(trip)
	return &Reverse{
		prx:  prx,
		repo: repo,
	}
}

type Reverse struct {
	prx  *httputil.ReverseProxy
	repo repository.All
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
	w, r := c.Response(), c.Request()
	ctx := r.Context()
	oid, _ := bson.ObjectIDFromHex(id)
	repo := rvs.repo.Agent()
	dat, err := repo.FindByID(ctx, oid)
	if err != nil {
		return err
	}
	brk := dat.Broker
	if brk == nil || brk.ID.IsZero() {
		return ship.ErrNotFound
	}

	// slash
	rawPath := r.URL.Path
	if pth != "/" && strings.HasSuffix(rawPath, "/") {
		pth += "/"
	}
	pth = "/api/reverse/agent/" + id + pth
	reqURL := transport.NewBrokerIDURL(brk.ID.Hex(), pth)
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
