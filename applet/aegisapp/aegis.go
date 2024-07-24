package aegisapp

import (
	"github.com/quic-go/quic-go"
)

type aegisApp struct{}

func (a *aegisApp) Proto() string {
	return "aegis"
}

func (a *aegisApp) Serve(conn quic.Connection) {
	// TODO implement me
	panic("implement me")
}
