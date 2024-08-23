package restapi

import (
	"context"
	"net/url"

	"github.com/coder/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/protocol/wsocket"
)

type forwardAPI struct{}

func (api *forwardAPI) Forward(c *ship.Context) error {
	if c.IsWebSocket() {
	}

	return nil
}

func (api *forwardAPI) ws() {
}

type forwardService struct {
	dest *url.URL
}

func (svc *forwardService) Websocket(cli *websocket.Conn, path string, quires url.Values) error {
	destURL := svc.URL(path, quires)
	srv, _, err := websocket.Dial(context.Background(), destURL.String(), nil)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer srv.CloseNow()

	wait := make(chan struct{})
	go func() {
		_ = wsocket.Copy(srv, cli)
		close(wait)
	}()
	_ = wsocket.Copy(cli, srv)
	<-wait

	return nil
}

func (svc *forwardService) URL(path string, quires url.Values) *url.URL {
	d := svc.dest
	u := (&url.URL{
		Scheme:      d.Scheme,
		Opaque:      d.Opaque,
		User:        d.User,
		Host:        d.Host,
		Path:        d.Path,
		RawPath:     d.RawPath,
		OmitHost:    d.OmitHost,
		ForceQuery:  d.ForceQuery,
		RawQuery:    d.RawQuery,
		Fragment:    d.Fragment,
		RawFragment: d.RawFragment,
	}).JoinPath(path)
	if len(quires) != 0 {
		query := u.Query()
		for k, vs := range quires {
			for _, v := range vs {
				query.Add(k, v)
			}
		}
		u.RawQuery = query.Encode()
	}

	return u
}
