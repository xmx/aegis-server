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

func (cl *connectListener) OnConnection(now time.Time, peer linkhub.Peer) {
}

func (cl *connectListener) OnDisconnection(now time.Time, info linkhub.Info) {
}
