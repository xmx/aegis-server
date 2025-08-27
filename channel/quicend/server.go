package quicend

import (
	"context"
	"net"

	"golang.org/x/net/quic"
)

type Server struct {
	end *quic.Endpoint
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) serve() error {
	conn, err := s.end.Accept(context.Background())
	if err != nil {
		return err
	}

	conn.Close()
}

// Listener
// Client
//
