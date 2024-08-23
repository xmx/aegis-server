package wsterm

import (
	"context"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/xmx/aegis-server/protocol/asciicast"
)

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{ws: ws}
}

type Conn struct {
	ws *websocket.Conn
}

func (c *Conn) Input() (code asciicast.CodeType, data string, err error) {
	lines := make([]string, 0, 2)
	if err = wsjson.Read(context.Background(), c.ws, &lines); err != nil {
		return "", "", err
	}
	for i, line := range lines {
		switch i {
		case 0:
			code = asciicast.CodeType(line)
		case 1:
			data = line
		}
	}

	return
}

func (c *Conn) Write(p []byte) (int, error) {
	n := len(p)
	ctx := context.Background()
	lines := []any{asciicast.CodeOutput, string(p)}
	if err := wsjson.Write(ctx, c.ws, lines); err != nil {
		return 0, err
	}

	return n, nil
}

func (c *Conn) Close() error {
	return c.ws.CloseNow()
}
