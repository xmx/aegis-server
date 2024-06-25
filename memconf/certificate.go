package memconf

import (
	"context"
	"crypto/tls"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/memoize"
)

type CertificateConfigurer interface {
	repository.CertificateRepository

	Certificate(context.Context) (credential.Certifier, error)
}

func Certificate(repo repository.CertificateRepository) CertificateConfigurer {
	c := &certificateConfig{repo: repo}
	c.cache = memoize.NewCache2(c.slowLoad)

	return c
}

type certificateConfig struct {
	repo  repository.CertificateRepository
	cache memoize.Cache2[credential.Certifier, error]
}

func (c *certificateConfig) Enabled(ctx context.Context) (*model.Certificate, error) {
	return c.repo.Enabled(ctx)
}

func (c *certificateConfig) Create(ctx context.Context, cert *model.Certificate) (bool, error) {
	return c.forget(c.repo.Create(ctx, cert))
}

func (c *certificateConfig) Delete(ctx context.Context, id int64) (bool, error) {
	return c.forget(c.repo.Delete(ctx, id))
}

func (c *certificateConfig) Certificate(ctx context.Context) (credential.Certifier, error) {
	return c.cache.Load(ctx)
}

func (c *certificateConfig) slowLoad(ctx context.Context) (credential.Certifier, error) {
	cert, err := c.Enabled(ctx)
	if err != nil {
		return nil, err
	}

	certPEMBlock := []byte(cert.PublicKey)
	keyPEMBlock := []byte(cert.PrivateKey)
	pair, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{pair}}
	pool := credential.Atomic()
	if err = pool.Modification(cfg); err != nil {
		return nil, err
	}

	return pool, nil
}

func (c *certificateConfig) forget(enabled bool, err error) (bool, error) {
	if err == nil && enabled {
		_, _ = c.cache.Forget()
	}

	return enabled, err
}
