package memconf

import (
	"context"
	"crypto/tls"
	"errors"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/memoize"
)

type ConfigCertificateConfigurer interface {
	repository.ConfigCertificateRepository
	credential.Certifier
}

func ConfigCertificate(repo repository.ConfigCertificateRepository) ConfigCertificateConfigurer {
	c := &configCertificateConfigurer{repo: repo}
	c.cache = memoize.NewCache2(c.slowLoad)

	return c
}

type configCertificateConfigurer struct {
	repo  repository.ConfigCertificateRepository
	cache memoize.Cache2[*tls.Config, error]
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

func (c *configCertificateConfigurer) Certificate(*tls.ClientHelloInfo) (*tls.Config, error) {
	return c.cache.Load(context.Background())
}

func (c *configCertificateConfigurer) Modification(*tls.Config) error {
	return errors.ErrUnsupported
}

func (c *configCertificateConfigurer) slowLoad(ctx context.Context) (*tls.Config, error) {
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

	return cfg, nil
}

func (c *configCertificateConfigurer) forget(enabled bool, err error) (bool, error) {
	if err == nil && enabled {
		_, _ = c.cache.Forget()
	}

	return enabled, err
}
