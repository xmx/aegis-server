package restapi

import (
	"io"
	"net"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/jsos/jsvm"
	"github.com/xmx/aegis-common/library/httpkit"
)

func NewPlay(mods []jsvm.Module) *Play {
	return &Play{
		mods: mods,
		wsu:  httpkit.NewWebsocketUpgrader(),
	}
}

type Play struct {
	mods []jsvm.Module
	wsu  *websocket.Upgrader
}

func (p *Play) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/play/js").GET(p.js)

	return nil
}

func (p *Play) js(c *ship.Context) error {
	//w, r := c.Response(), c.Request()
	//ws, err := p.wsu.Upgrade(w, r, nil)
	//if err != nil {
	//	return err
	//}
	//defer ws.Close()
	//
	//wsout := &playWriter{channel: "stdout", socket: ws}
	//wserr := &playWriter{channel: "stderr", socket: ws}
	//req := new(request.PlayJS)
	//now := time.Now()
	//_ = ws.SetReadDeadline(now.Add(10 * time.Second))
	//if err = ws.ReadJSON(req); err != nil {
	//	p.writeError(wserr, err)
	//	return err
	//}
	//if err = c.Validator.Validate(req); err != nil {
	//	p.writeError(wserr, err)
	//	return err
	//}
	//ch, msg := req.Channel, req.Message
	//if ch != "stdin" {
	//	p.writeText(wserr, "当前仅允许 channel 为 stdin 的消息")
	//	return nil
	//}
	//
	//mods := append(p.mods, jsstd.All()...)
	//vm := jsvm.NewVM(context.Background(), slog.Default())
	//require := vm.Require()
	//require.Registers(mods)
	//stdout, stderr := vm.Output()
	//stdout.Attach(wsout)
	//stderr.Attach(wserr)
	//defer func() {
	//	stdout.Detach(wsout)
	//	stderr.Detach(wserr)
	//}()
	//
	//wait := make(chan struct{})
	//go func() {
	//	defer func() {
	//		_ = ws.Close()
	//		vm.Kill(net.ErrClosed)
	//		close(wait)
	//	}()
	//	val, exx := vm.RunScript(ch, msg)
	//	if exx != nil {
	//		_, _ = wsout.Write([]byte(exx.Error()))
	//	} else if val != nil && sobek.IsUndefined(val) {
	//		_, _ = wserr.Write([]byte(val.String()))
	//	}
	//}()
	//p.read(ws, vm)
	//
	//<-wait

	return nil
}

func (p *Play) read(ws *websocket.Conn, vm jsvm.Engineer) {
	defer vm.Kill(net.ErrClosed)
	for {
		_, r, err := ws.NextReader()
		if err != nil {
			break
		}
		_, _ = io.Copy(io.Discard, r)
	}
}

func (p *Play) writeError(w io.Writer, err error) {
	var msg string
	switch err.(type) {
	case nil:
		msg = "<no error>"
	case interface{ Timeout() bool }:
		msg = "消息超时"
	default:
		msg = err.Error()
	}

	p.writeText(w, msg)
}

func (p *Play) writeText(w io.Writer, msg string) {
	_, _ = w.Write([]byte(msg))
}

type playWriter struct {
	channel string
	socket  *websocket.Conn
}

func (pw *playWriter) Write(p []byte) (int, error) {
	n := len(p)
	data := &playData{Channel: pw.channel, Message: string(p)}
	if err := pw.socket.WriteJSON(data); err != nil {
		return 0, err
	}

	return n, nil
}

type playData struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}
