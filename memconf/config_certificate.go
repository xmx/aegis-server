package memconf

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/memoize"
)

type ConfigCertificateConfigurer interface {
	repository.ConfigCertificateRepository
	Certificate(hi *tls.ClientHelloInfo) (*tls.Config, error)
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

func (c *configCertificateConfigurer) Update(ctx context.Context, cert *model.ConfigCertificate) (bool, error) {
	return c.forget(c.repo.Update(ctx, cert))
}

func (c *configCertificateConfigurer) Delete(ctx context.Context, id int64) (bool, error) {
	return c.forget(c.repo.Delete(ctx, id))
}

func (c *configCertificateConfigurer) Certificate(hi *tls.ClientHelloInfo) (*tls.Config, error) {
	ctx := hi.Context()
	return c.cache.Load(ctx)
}

func (c *configCertificateConfigurer) slowLoad(ctx context.Context) (*tls.Config, error) {
	cert, err := c.Enabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("服务端查询证书错误: %w", err)
	}

	certPEMBlock := []byte(cert.PublicKey)
	keyPEMBlock := []byte(cert.PrivateKey)
	pair, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, fmt.Errorf("服务端证书格式错误: %w", err)
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
