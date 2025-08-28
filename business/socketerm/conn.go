package socketerm

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

func New(ws *websocket.Conn, timeout time.Duration) *Conn {
	return &Conn{
		ws:      ws,
		timeout: timeout,
	}
}

type Conn struct {
	ws      *websocket.Conn
	timeout time.Duration
}

func (c *Conn) Write(p []byte) (int, error) {
	lines := []any{"o", string(p)}
	if err := c.ws.WriteJSON(lines); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (c *Conn) Recv() (kind, data string, err error) {
	var lines []string
	if err = c.ws.ReadJSON(&lines); err != nil {
		return
	}
	for i, line := range lines {
		switch i {
		case 0:
			kind = line
		case 1:
			data = line
		}
	}

	return
}

func (c *Conn) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}
