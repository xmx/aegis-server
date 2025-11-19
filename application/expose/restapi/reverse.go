package restapi

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/library/httpkit"
	"github.com/xmx/aegis-common/tunnel/tunconst"
	"github.com/xmx/aegis-common/tunnel/tundial"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/library/httpnet"
)

func NewReverse(dial tundial.ContextDialer, repo repository.All) *Reverse {
	trip := &http.Transport{DialContext: dial.DialContext}
	prx := httpnet.NewReverse(trip)
	wsd := httpkit.NewWebsocketDialer(dial.DialContext)

	return &Reverse{
		prx:  prx,
		wsd:  wsd,
		wsu:  httpkit.NewWebsocketUpgrader(),
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

	reqURL := tunconst.ServerToAgent(id, "/api/reverse/")
	reqURL = reqURL.JoinPath(id, pth)
	reqURL.RawQuery = r.URL.RawQuery
	r.URL = reqURL
	r.Host = reqURL.Host

	rvs.prx.ServeHTTP(w, r)

	return nil
}

func (rvs *Reverse) broker(c *ship.Context) error {
	id, path := c.Param("id"), "/"+c.Param("path")
	w, r := c.Response(), c.Request()
	reqURL := r.URL
	reqPath := reqURL.Path
	if path != "/" && strings.HasSuffix(reqPath, "/") {
		path += "/"
	}

	destURL := tunconst.ServerToBroker(id, path)
	destURL.RawQuery = reqURL.RawQuery

	if c.IsWebSocket() {
		rvs.serveWebsocket(c, destURL)
		return nil
	}

	r.URL = destURL
	r.Host = destURL.Host
	rvs.prx.ServeHTTP(w, r)

	return nil
}

func (rvs *Reverse) serveWebsocket(c *ship.Context, destURL *url.URL) {
	w, r := c.Response(), c.Request()
	ctx := r.Context()

	cli, err := rvs.wsu.Upgrade(w, r, nil)
	if err != nil {
		c.Errorf("websocket upgrade 失败", "error", err)
		return
	}
	defer cli.Close()

	destURL.Scheme = "ws"
	srv, _, err := rvs.wsd.DialContext(ctx, destURL.String(), nil)
	if err != nil {
		c.Errorf("websocket 后端连接失败", "error", err)
		_ = rvs.writeClose(cli, err)
		return
	}
	defer srv.Close()

	ret := httpkit.ExchangeWebsocket(cli, srv)
	c.Infof("websocket 连接结束", slog.Any("result", ret))
}

func (rvs *Reverse) writeClose(cli *websocket.Conn, err error) error {
	return cli.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
}
