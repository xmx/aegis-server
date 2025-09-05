package restapi

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/contract/problem"
	"github.com/xmx/aegis-server/applet/expose/service"
	"github.com/xmx/aegis-server/channel/broker"
	"github.com/xmx/aegis-server/contract/request"
)

func NewBroker(svc *service.Broker, trp http.RoundTripper) *Broker {
	prx := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetXForwarded()
		},
		Transport: trp,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			code := http.StatusBadGateway
			pb := &problem.Details{
				Host:     r.Host,
				Type:     r.URL.Host,
				Status:   code,
				Detail:   err.Error(),
				Instance: r.URL.Path,
				Method:   r.Method,
				Datetime: time.Now().UTC(),
			}
			if ae, ok := err.(*net.AddrError); ok {
				pb.Detail = ae.Err
			}
			_ = pb.JSON(w)
		},
	}

	return &Broker{
		svc: svc,
		prx: prx,
	}
}

type Broker struct {
	svc *service.Broker
	prx *httputil.ReverseProxy
}

func (bk *Broker) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/brokers").GET(bk.list)
	r.Route("/broker").POST(bk.create)
	r.Route("/broker/kickout").GET(bk.kickout)
	r.Route("/broker/reverse/:id/").Any(bk.reverse)
	r.Route("/broker/reverse/:id/*path").Any(bk.reverse)
	return nil
}

func (bk *Broker) create(c *ship.Context) error {
	req := new(request.BrokerCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return bk.svc.Create(ctx, req.Name)
}

func (bk *Broker) list(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := bk.svc.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (bk *Broker) kickout(c *ship.Context) error {
	req := new(request.ObjectID)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	return bk.svc.Kickout(req.OID())
}

func (bk *Broker) reverse(c *ship.Context) error {
	id, path := c.Param("id"), "/"+c.Param("path")
	w, r := c.Response(), c.Request()
	reqURL := broker.MakesBrokerURL(id, path)
	r.URL = reqURL
	r.Host = reqURL.Host

	bk.prx.ServeHTTP(w, r)

	return nil
}
