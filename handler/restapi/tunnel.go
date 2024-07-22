package restapi

import (
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/xgfone/ship/v5"
)

func DD(c *ship.Context) error {
	w := c.ResponseWriter()
	hc := w.(http3.Hijacker).Connection()
	conn := hc.(quic.Connection)

	trip := &http3.SingleDestinationRoundTripper{
		Connection: conn,
	}
	cli := &http.Client{
		Transport: trip,
	}
	_ = cli

	return nil
}
