package restapi

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/handler/shipx"
)

func REPL() shipx.Register {
	upgrade := &websocket.Upgrader{
		HandshakeTimeout:  10 * time.Second,
		CheckOrigin:       func(*http.Request) bool { return true },
		EnableCompression: true,
	}

	return &replAPI{
		upgrade: upgrade,
	}
}

type replAPI struct {
	upgrade *websocket.Upgrader
}

func (api *replAPI) Register(rt shipx.Router) error {
	rt.Auth().Route("/repl/js").GET(api.JS)
	return nil
}

func (api *replAPI) JS(c *ship.Context) error {
	ws, err := api.upgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer ws.Close()

	return nil
}
