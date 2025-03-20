package socketerm

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
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
	ctx, cancel := c.getContext()
	defer cancel()

	lines := []any{"o", string(p)}
	if err := wsjson.Write(ctx, c.ws, lines); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (c *Conn) Recv() (kind, data string, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	lines := make([]string, 0, 2)
	if err = wsjson.Read(ctx, c.ws, &lines); err != nil {
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
