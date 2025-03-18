package restapi

import (
	"context"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/grafana/sobek"
	"github.com/xmx/aegis-server/jsenv/jsmod"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/ship"
)

func NewPlay(mods []jsvm.GlobalRegister) *Play {
	return &Play{
		mods: mods,
	}
}

type Play struct {
	mods []jsvm.GlobalRegister
}

func (ply *Play) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/ws/play/js").GET(ply.run)
	return nil
}

func (ply *Play) run(c *ship.Context) error {
	w, r := c.ResponseWriter(), c.Request()
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	vm := jsvm.New()
	out := &wsout{ws: ws}
	mods := append(ply.mods, jsmod.NewConsole(out))
	if err = jsvm.RegisterGlobals(vm, mods); err != nil {
		return err
	}

	ctx := r.Context()
	for {
		data := new(Data)
		if err = wsjson.Read(ctx, ws, data); err != nil {
			break
		}
		if val, err := vm.RunString(data.Code); err != nil {
			out.WriteError(err)
		} else if val != nil && val != sobek.Undefined() {
			out.Write([]byte(val.String()))
		}
	}

	return nil
}

type Data struct {
	Code string `json:"code"`
}

type Response struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type wsout struct {
	ws *websocket.Conn
}

func (w *wsout) Write(p []byte) (int, error) {
	data := &Response{Type: "stdout", Data: string(p)}
	err := wsjson.Write(context.Background(), w.ws, data)
	return len(p), err
}

func (w *wsout) WriteError(err error) {
	data := &Response{Type: "stderr", Data: err.Error()}
	wsjson.Write(context.Background(), w.ws, data)
}
