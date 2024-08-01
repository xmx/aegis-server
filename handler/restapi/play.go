package restapi

import (
	"context"
	"errors"
	"mime"
	"net"
	"net/http"
	"time"

	"github.com/dop251/goja"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/jsenv/jslib"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/protocol/wsocket"
	"nhooyr.io/websocket"
)

func NewPlay(player service.Player) shipx.Register {
	return &playAPI{
		player: player,
	}
}

type playAPI struct {
	player service.Player
}

func (api *playAPI) Register(rt shipx.Router) error {
	auth := rt.Auth()
	auth.Route("/play/js").GET(api.JS)
	auth.Route("/play/pprof").GET(api.Pprof)
	return nil
}

func (api *playAPI) JS(c *ship.Context) error {
	w, r := c.ResponseWriter(), c.Request()
	opt := &websocket.AcceptOptions{
		CompressionMode:    websocket.CompressionContextTakeover,
		InsecureSkipVerify: true,
	}
	ws, err := websocket.Accept(w, r, opt)
	if err != nil {
		return err
	}
	conn := wsocket.NewConn(ws)
	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()

	stdout := wsocket.JSWriter(conn, wsocket.KindStdout) // 重定向输出数据
	player := api.player.NewGoja([]jsvm.Loader{jslib.Console(stdout)})

	valid := c.Validator
	for {
		req := new(request.PlayJS)
		if err = conn.ReadJSON(req); err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			_, _, detail := shipx.UnpackError(err)
			_ = conn.WriteJSON(wsocket.SErrorBody(detail))
			continue
		}
		if err = valid.Validate(req); err != nil {
			_, _, detail := shipx.UnpackError(err)
			_ = conn.WriteJSON(wsocket.SErrorBody(detail))
			continue
		}

		args := []jsvm.Loader{jslib.ArgsPrototype(req.Args)}
		_, err = player.Exec(context.Background(), args, req.Script)
		if err != nil {
			c.Infof("运行脚本错误： %v", err)
			_ = conn.WriteJSON(wsocket.ErrorBody(err))
		}
	}

	return nil
}

func (api *playAPI) Pprof(c *ship.Context) error {
	w := c.ResponseWriter()
	w.Header().Set(ship.HeaderContentType, ship.MIMEOctetStream)

	now := time.Now()
	name := now.Format(time.RFC3339)
	param := map[string]string{"filename": name + ".pprof"}
	disposition := mime.FormatMediaType("attachment", param)
	w.Header().Set(ship.HeaderContentDisposition, disposition)
	w.WriteHeader(http.StatusOK)

	if err := goja.StartProfile(w); err != nil {
		return err
	}
	defer goja.StopProfile()

	ctx := c.Request().Context()
	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}

	return nil
}
