package memconf

import (
	"context"
	"crypto/tls"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/memoize"
)

type ConfigCertificateConfigurer interface {
	repository.ConfigCertificateRepository
	credential.Certifier
	// Certificate(context.Context) (credential.Certifier, error)
}

func ConfigCertificate(repo repository.ConfigCertificateRepository) ConfigCertificateConfigurer {
	c := &configCertificateConfigurer{repo: repo}
	c.cache = memoize.NewCache2(c.slowLoad)

	return c
}

type configCertificateConfigurer struct {
	repo  repository.ConfigCertificateRepository
	cache memoize.Cache2[credential.Certifier, error]
}

func (c *configCertificateConfigurer) Enabled(ctx context.Context) (*model.ConfigCertificate, error) {
	return c.repo.Enabled(ctx)
}

func (c *configCertificateConfigurer) Create(ctx context.Context, cert *model.ConfigCertificate) (bool, error) {
	return c.forget(c.repo.Create(ctx, cert))
}

func (c *configCertificateConfigurer) Delete(ctx context.Context, id int64) (bool, error) {
	return c.forget(c.repo.Delete(ctx, id))
}

func (c *configCertificateConfigurer) Certificate(ctx context.Context) (credential.Certifier, error) {
	return c.cache.Load(ctx)
}

func (c *configCertificateConfigurer) slowLoad(ctx context.Context) (credential.Certifier, error) {
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

func (c *configCertificateConfigurer) forget(enabled bool, err error) (bool, error) {
	if err == nil && enabled {
		_, _ = c.cache.Forget()
	}

	return enabled, err
}
