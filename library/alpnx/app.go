package alpnx

import (
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type Application interface {
	Proto() string
	Serve(conn quic.Connection) error
}

func HTTP3(srv *http3.Server) Application {
	return &http3App{srv: srv}
}

type http3App struct {
	srv *http3.Server
}

func (h *http3App) Proto() string {
	return "h3"
}

func (h *http3App) Serve(conn quic.Connection) error {
	return h.srv.ServeQUICConn(conn)
}
