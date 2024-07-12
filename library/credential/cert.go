package credential

import (
	"crypto/tls"
	"maps"
	"strings"
	"sync/atomic"
)

type Certifier interface {
	Get(*tls.ClientHelloInfo) (*tls.Certificate, error)
	Put(name string, cert *tls.Certificate)
	Clear()
}

func NewPool() Certifier {
	return &certPool{}
}

type certPool struct {
	base *tls.Config
	val  atomic.Value
}

func (a *certPool) Get(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	hm, yes := a.val.Load().(map[string]*tls.Certificate)
	if !yes || len(hm) == 0 {
		return nil, nil
	}

	// https://github.com/golang/go/blob/go1.22.5/src/crypto/tls/common.go#L1141-L1154
	name := strings.ToLower(info.ServerName)
	if cert, ok := hm[name]; ok && cert != nil || name == "" {
		return cert, nil
	}

	labels := strings.Split(name, ".")
	labels[0] = "*"
	wildcardName := strings.Join(labels, ".")
	cert := hm[wildcardName]

	return cert, nil
}

// Put 存放证书。
func (a *certPool) Put(name string, cert *tls.Certificate) {
	hm, ok := a.val.Load().(map[string]*tls.Certificate)
	if !ok {
		a.val.Store(map[string]*tls.Certificate{name: cert})
		return
	}

	certs := maps.Clone(hm)
	certs[name] = cert
	a.val.Store(certs)
}

func (a *certPool) Clear() {
	a.val.Store(map[string]*tls.Certificate{})
}
