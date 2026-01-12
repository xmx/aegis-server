package restapi

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/muxlink/muxproto"
	"github.com/xmx/aegis-common/problem"
	"github.com/xmx/aegis-common/wsocket"
	"github.com/xmx/aegis-server/application/errcode"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewReverse(dial muxproto.Dialer) *Reverse {
	trip := &http.Transport{DialContext: dial.DialContext}
	resv := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetXForwarded()
		},
		Transport: trip,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			host := r.Host
			if host == "" {
				host = r.URL.Host
			}

			prob := &problem.Details{
				Host:     host,
				Status:   http.StatusBadGateway,
				Instance: r.URL.Path,
				Method:   r.Method,
				Datetime: time.Now().UTC(),
			}

			if ae, ok := err.(*net.OpError); ok {
				addr := ae.Addr.String()
				if brokID, _, found := strings.Cut(addr, muxproto.BrokerHostSuffix); found {
					err = errcode.FmtBrokerDisconnect.Fmt(brokID)
				} else if errors.Is(ae.Err, mongo.ErrNoDocuments) {
					if agentID, _, exists := strings.Cut(addr, muxproto.AgentHostSuffix); exists {
						err = errcode.FmtAgentNotExists.Fmt(agentID)
					}
				}
			}
			prob.Detail = err.Error()

			_ = prob.JSON(w)
		},
	}
	wsd := &websocket.Dialer{
		NetDialContext:   dial.DialContext,
		HandshakeTimeout: 10 * time.Second,
	}
	wsu := &websocket.Upgrader{
		HandshakeTimeout:  10 * time.Second,
		CheckOrigin:       func(*http.Request) bool { return true },
		EnableCompression: true,
	}

	return &Reverse{
		prx: resv,
		wsd: wsd,
		wsu: wsu,
	}
}

type Reverse struct {
	prx *httputil.ReverseProxy
	wsd *websocket.Dialer
	wsu *websocket.Upgrader
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

	reqURL := muxproto.ServerToAgentURL(id, "/api/reverse/")
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

	destURL := muxproto.ServerToBrokerURL(id, path)
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

	ret := wsocket.Exchange(cli, srv)
	c.Infof("websocket 连接结束", slog.Any("result", ret))
}

func (rvs *Reverse) writeClose(cli *websocket.Conn, err error) error {
	return cli.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
}
