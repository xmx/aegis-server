package tunutil

import (
	"log/slog"
	"time"

	"github.com/xmx/aegis-control/linkhub"
)

func NewConnectListener(log *slog.Logger) linkhub.ConnectListener {
	return &connectListener{log: log}
}

type connectListener struct {
	log *slog.Logger
}

func (cl *connectListener) OnConnection(peer linkhub.Peer, connectAt time.Time) {
}

func (cl *connectListener) OnDisconnection(info linkhub.Info, connectAt, disconnectAt time.Time) {
}
