package restapi

import (
	"context"
	"io"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/dop251/goja"
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

	sig := new(singleVM)
	ctx := r.Context()
	context.AfterFunc(ctx, sig.kill)
	type tempData struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}
	for {
		data := new(tempData)
		if err = wsjson.Read(ctx, ws, data); err != nil {
			break
		}

		switch data.Type {
		case "kill":
			sig.kill()
		case "exec":
			vm := jsvm.New()
			sig.set(vm)
			go ply.newInstanceExec(vm, ws, data.Data)
		}
	}

	return nil
}

func (ply *Play) newInstanceExec(vm jsvm.Runtime, ws *websocket.Conn, code string) error {
	stdout, stderr := ply.stdout(ws), ply.stderr(ws)
	mods := append(ply.mods, jsmod.NewConsole(stdout))
	if err := vm.RegisterGlobals(mods); err != nil {
		return err
	}

	val := vm.Runtime().Get("os")
	if obj, _ := val.(*goja.Object); obj != nil {
		_ = obj.Set("stdout", stdout)
		_ = obj.Set("stderr", stderr)
	}
	if ret, exx := vm.RunString(code); exx != nil {
		_, _ = stderr.Write([]byte(exx.Error()))
		return exx
	} else if ret != nil && ret != goja.Undefined() {
		_, _ = stdout.Write([]byte(ret.String()))
	}

	return nil
}

func (ply *Play) stderr(ws *websocket.Conn) io.Writer {
	return &socketConn{tp: "stderr", ws: ws}
}

func (ply *Play) stdout(ws *websocket.Conn) io.Writer {
	return &socketConn{tp: "stdout", ws: ws}
}

type socketConn struct {
	tp string
	ws *websocket.Conn
}

func (sc *socketConn) Write(p []byte) (int, error) {
	data := struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}{
		Type: sc.tp,
		Data: string(p),
	}
	err := wsjson.Write(context.Background(), sc.ws, data)

	return len(p), err
}

type singleVM struct {
	mu sync.Mutex
	vm jsvm.Runtime
}

func (sig *singleVM) set(vm jsvm.Runtime) {
	sig.mu.Lock()
	sig.interrupt()
	sig.vm = vm
	sig.mu.Unlock()
}

func (sig *singleVM) kill() {
	sig.mu.Lock()
	sig.interrupt()
	sig.mu.Unlock()
}

func (sig *singleVM) interrupt() {
	if vm := sig.vm; vm != nil {
		vm.Interrupt("killed")
	}
}
