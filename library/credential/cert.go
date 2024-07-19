package credential

import (
	"crypto/tls"
	"strings"
	"sync/atomic"
)

type Certifier interface {
	Match(*tls.ClientHelloInfo) (*tls.Config, error)
	Replace(certs []tls.Certificate)
}

func Pool(base *tls.Config) Certifier {
	return &certPool{base: base}
}

type certPool struct {
	base *tls.Config
	val  atomic.Value
}

func (cp *certPool) Match(info *tls.ClientHelloInfo) (*tls.Config, error) {
	hm, yes := cp.val.Load().(map[string]*tls.Config)
	if !yes || len(hm) == 0 {
		return nil, nil
	}

	// https://github.com/golang/go/blob/go1.22.5/src/crypto/tls/common.go#L1141-L1154
	name := strings.ToLower(info.ServerName)
	if cfg, ok := hm[name]; ok && cfg != nil || name == "" {
		return cfg, nil
	}

	labels := strings.Split(name, ".")
	labels[0] = "*"
	wildcardName := strings.Join(labels, ".")
	cfg := hm[wildcardName]

	return cfg, nil
}

func (cp *certPool) Replace(certs []tls.Certificate) {
	hm := make(map[string]*tls.Config, len(certs)*4)
	for _, cert := range certs {
		leaf := cert.Leaf
		for _, name := range leaf.DNSNames {
			cfg := cp.base.Clone()
			cfg.Certificates = []tls.Certificate{cert}
			hm[name] = cfg
		}
		for _, ip := range leaf.IPAddresses {
			cfg := cp.base.Clone()
			cfg.Certificates = []tls.Certificate{cert}
			hm[ip.String()] = cfg
		}
	}

	cp.val.Store(hm)
}
