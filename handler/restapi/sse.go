package restapi

import (
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/protocol/eventsource"
)

type sseAPI struct{}

func (api *sseAPI) Register(rt shipx.Router) error {
	// TODO implement me
	panic("implement me")
}

func (api *sseAPI) Multi(c *ship.Context) error {
	c.IsWebSocket()
	w, r := c.ResponseWriter(), c.Request()
	sse := eventsource.Accept(w, r)
	if sse == nil {
		return errcode.ErrServerSentEvents
	}
	//goland:noinspection GoUnhandledErrorResult
	defer sse.Close()

	<-sse.Done()

	return nil
}
