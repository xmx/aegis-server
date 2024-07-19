package restapi

import (
	"github.com/quic-go/quic-go/http3"
	"github.com/xgfone/ship/v5"
)

func DD(c *ship.Context) error {
	body := c.Body()
	streamer := body.(http3.HTTPStreamer)
	stm := streamer.HTTPStream()
	stm.StreamID()

	return nil
}
