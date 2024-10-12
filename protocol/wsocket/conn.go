package wsocket

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func KindConn(ws *websocket.Conn, kind Kind) *Conn {
	return &Conn{ws: ws, kind: kind}
}

type Conn struct {
	ws   *websocket.Conn
	kind Kind
}

func (c *Conn) ReadJSON(v any) error {
	return wsjson.Read(context.Background(), c.ws, v)
}

func (c *Conn) Recv() (*Recv, error) {
	v := new(Recv)
	if err := c.ReadJSON(v); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *Conn) WriteJSON(v any) error {
	return wsjson.Write(context.Background(), c.ws, v)
}

func (c *Conn) Write(p []byte) (int, error) {
	n := len(p)
	v := &Body{Kind: c.kind, Data: string(p)}
	if err := c.WriteJSON(v); err != nil {
		return 0, err
	}

	return n, nil
}

func (c *Conn) Close() error {
	return c.ws.CloseNow()
}

func (c *Conn) CloseJSON(code websocket.StatusCode, v any) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	return c.ws.Close(code, buf.String())
}

type ConnWriter struct {
	ws *websocket.Conn
}
