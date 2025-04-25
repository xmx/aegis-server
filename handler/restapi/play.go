package restapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/jsos/jsvm"
	"github.com/xmx/ship"
)

func NewPlay(mods []jsvm.ModuleRegister) *Play {
	return &Play{
		mods: mods,
		dir:  filepath.Join("resources", "app"),
	}
}

type Play struct {
	mods []jsvm.ModuleRegister
	dir  string
}

func (ply *Play) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/play/upload").PUT(ply.upload)
	r.Route("/ws/play/run").GET(ply.run)
	return nil
}

func (ply *Play) upload(c *ship.Context) error {
	req := new(request.PlayUpload)
	if err := c.Bind(req); err != nil {
		return err
	}
	file, err := req.File.Open()
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	tmp := make([]byte, 16)
	_, _ = rand.Read(tmp)
	uid := hex.EncodeToString(tmp)
	name := uid + ".zip"
	fp := filepath.Join(ply.dir, name)

	dest, err := os.Create(fp)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		return err
	}
	ret := &response.PlayUpload{ID: uid}

	return c.JSON(http.StatusOK, ret)
}

func (ply *Play) run(c *ship.Context) error {
	uid := c.Query("id")
	name := uid + ".zip"
	fp := filepath.Join(ply.dir, filepath.Clean(name))

	eng, err := jsvm.New(ply.mods...)
	if err != nil {
		return err
	}
	defer eng.Kill("程序执行完毕")

	w, r := c.ResponseWriter(), c.Request()
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer ws.CloseNow()

	stderr := ply.stderr(ws)
	eng.Device().Stdout().Attach(ply.stdout(ws))
	eng.Device().Stderr().Attach(stderr)

	parent := r.Context()
	ctx, cancel := context.WithCancel(parent)

	go func() {
		if _, err = eng.RunJZip(fp); err != nil {
			_, _ = stderr.Write([]byte(err.Error()))
		}
		cancel()
	}()

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
			eng.Kill("客户端 killed")
		}
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
