package tunutil

import (
	"log/slog"
	"time"

	"github.com/xmx/aegis-control/linkhub"
)

func NewConnectListener(log *slog.Logger) linkhub.ServerHooker {
	return &connectListener{log: log}
}

type connectListener struct {
	log *slog.Logger
}

func (cl *connectListener) OnConnected(info linkhub.Info, connectAt time.Time) {
}

func (cl *connectListener) OnDisconnected(info linkhub.Info, connectAt, disconnectAt time.Time) {
}
