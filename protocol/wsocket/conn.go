package wsocket

import (
	"bytes"
	"context"
	"encoding/json"

	"nhooyr.io/websocket"
)

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{ws: ws}
}

type Conn struct {
	ws *websocket.Conn
}

func (c *Conn) ReadJSON(v any) error {
	_, reader, err := c.ws.Reader(context.Background())
	if err != nil {
		return err
	}

	return json.NewDecoder(reader).Decode(v)
}

func (c *Conn) WriteJSON(v any) error {
	writer, err := c.ws.Writer(context.Background(), websocket.MessageText)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer writer.Close()

	return json.NewEncoder(writer).Encode(v)
}

func (c *Conn) Close() error {
	return c.ws.CloseNow()
}

func (c *Conn) CloseJSON(code websocket.StatusCode, reason any) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(reason); err != nil {
		return err
	}

	return c.ws.Close(code, buf.String())
}
