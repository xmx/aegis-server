package restapi

import (
	"context"
	"errors"
	"mime"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/dop251/goja"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/jsenv/jslib"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/protocol/wsocket"
)

func NewPlay(player service.Player) *Play {
	return &Play{
		player: player,
	}
}

type Play struct {
	player service.Player
}

func (api *Play) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/ws/js/play").GET(api.js)
	r.Route("/js/play/pprof").GET(api.pprof)
	return nil
}

func (api *Play) js(c *ship.Context) error {
	w, r := c.ResponseWriter(), c.Request()
	opt := &websocket.AcceptOptions{
		CompressionMode:    websocket.CompressionContextTakeover,
		InsecureSkipVerify: true,
	}
	ws, err := websocket.Accept(w, r, opt)
	if err != nil {
		return err
	}
	conn := wsocket.KindConn(ws, wsocket.KindStdout)
	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()

	player := api.player.NewGoja([]jsvm.Loader{jslib.Console(conn)})
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

func (api *Play) pprof(c *ship.Context) error {
	secVal := c.Query("seconds")
	sec, _ := strconv.ParseInt(secVal, 10, 64)
	if sec < 1 {
		sec = 10
	}

	w, r := c.ResponseWriter(), c.Request()
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

	ctx := r.Context()
	timer := time.NewTimer(time.Duration(sec) * time.Second)

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
	goja.StopProfile()
	timer.Stop()

	return nil
}
