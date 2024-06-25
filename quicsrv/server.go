package quicsrv

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/xmx/aegis-server/library/credential"
)

type Server struct {
	Cert               credential.Certifier
	Addr               string
	Port               int
	QUICConfig         *quic.Config
	Handler            http.Handler
	EnableDatagrams    bool
	MaxHeaderBytes     int
	AdditionalSettings map[uint64]uint64
	StreamHijacker     func(http3.FrameType, quic.ConnectionTracingID, quic.Stream, error) (hijacked bool, err error)
	UniStreamHijacker  func(http3.StreamType, quic.ConnectionTracingID, quic.ReceiveStream, error) (hijacked bool)
	ConnContext        func(ctx context.Context, c quic.Connection) context.Context
	Logger             *slog.Logger
	server             atomic.Value
}

func (s *Server) ListenAndServe() error {
	tlsCfg := &tls.Config{
		GetConfigForClient: s.Cert.Certificate,
	}
	srv := &http3.Server{
		Addr:               s.Addr,
		Port:               s.Port,
		TLSConfig:          tlsCfg,
		QUICConfig:         s.QUICConfig,
		Handler:            s.Handler,
		EnableDatagrams:    s.EnableDatagrams,
		MaxHeaderBytes:     s.MaxHeaderBytes,
		AdditionalSettings: s.AdditionalSettings,
		StreamHijacker:     s.StreamHijacker,
		UniStreamHijacker:  s.UniStreamHijacker,
		ConnContext:        s.ConnContext,
		Logger:             s.Logger,
	}
	s.server.Store(srv)

	return srv.ListenAndServe()
}

func (s *Server) Close() error {
	if srv, _ := s.server.Load().(*http3.Server); srv != nil {
		return srv.Close()
	}

	return nil
}
