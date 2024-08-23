package wsocket

import (
	"context"

	"github.com/coder/websocket"
)

func Copy(dest, src *websocket.Conn) error {
	for {
		mt, data, err := src.Read(context.Background())
		if err != nil {
			_ = dest.CloseNow()
			return err
		}
		if err = dest.Write(context.Background(), mt, data); err != nil {
			_ = src.CloseNow()
			return err
		}
	}
}
