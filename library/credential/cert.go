package credential

import (
	"crypto/tls"
	"sync/atomic"
)

type Certifier interface {
	Certificate(*tls.ClientHelloInfo) (*tls.Config, error)
	Modification(cfg *tls.Config) error
}

func Atomic() Certifier {
	return new(singleCert)
}

type singleCert struct {
	val atomic.Value
}

func (s *singleCert) Certificate(*tls.ClientHelloInfo) (*tls.Config, error) {
	cfg, _ := s.val.Load().(*tls.Config)
	return cfg, nil
}

func (s *singleCert) Modification(cfg *tls.Config) error {
	s.val.Store(cfg)
	return nil
}
