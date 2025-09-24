package restapi

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/library/wsocket"
	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/library/httpnet"
)

func NewReverse(dial transport.Dialer, repo repository.All) *Reverse {
	trip := &http.Transport{DialContext: dial.DialContext}
	prx := httpnet.NewReverse(trip)
	wsd := &websocket.Dialer{
		NetDialContext:   dial.DialContext,
		HandshakeTimeout: 5 * time.Second,
	}

	return &Reverse{
		prx:  prx,
		wsd:  wsd,
		wsu:  wsocket.NewUpgrade(),
		repo: repo,
	}
}

type Reverse struct {
	prx  *httputil.ReverseProxy
	wsd  *websocket.Dialer
	wsu  *websocket.Upgrader
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

	rawPath := r.URL.Path
	if pth != "/" && strings.HasSuffix(rawPath, "/") {
		pth += "/"
	}
	pth = "/api/reverse/agent/" + id + pth
	reqURL := transport.NewServerBrokerAgentURL(id, pth)
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

	reqURL := transport.NewServerBrokerURL(id, path)
	r.URL = reqURL
	r.Host = reqURL.Host

	if c.IsWebSocket() {
		rvs.serveWebsocket(c, reqURL)
	} else {
		rvs.prx.ServeHTTP(w, r)
	}

	return nil
}

func (rvs *Reverse) serveWebsocket(c *ship.Context, reqURL *url.URL) {
	w, r := c.Response(), c.Request()
	ctx := r.Context()

	cli, err := rvs.wsu.Upgrade(w, r, nil)
	if err != nil {
		c.Errorf("websocket upgrade 失败", "error", err)
		return
	}
	defer cli.Close()

	switch reqURL.Scheme {
	case "https":
		reqURL.Scheme = "wss"
	default:
		reqURL.Scheme = "ws"
	}
	srv, _, err := rvs.wsd.DialContext(ctx, reqURL.String(), nil)
	if err != nil {
		c.Errorf("websocket 后端连接失败", "error", err)
		_ = rvs.writeClose(cli, err)
		return
	}
	defer srv.Close()

	ret := wsocket.Pipe(cli, srv)
	c.Infof("websocket 连接结束", slog.Any("result", ret))
}

func (rvs *Reverse) writeClose(cli *websocket.Conn, err error) error {
	return cli.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
}
