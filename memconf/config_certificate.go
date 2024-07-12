package memconf

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/mcache"
)

type ConfigCertificateConfigurer interface {
	repository.ConfigCertificateRepository
	Certificate(*tls.ClientHelloInfo) (*tls.Certificate, error)
}

func ConfigCertificate(repo repository.ConfigCertificateRepository) ConfigCertificateConfigurer {
	c := &configCertificateConfigurer{repo: repo}
	c.cache = mcache.NewCache2(c.slowLoad)

	return c
}

type configCertificateConfigurer struct {
	repo  repository.ConfigCertificateRepository
	cache mcache.Cache2[credential.Certifier, error]
}

func (c *configCertificateConfigurer) Enables(ctx context.Context) ([]*model.ConfigCertificate, error) {
	return c.repo.Enables(ctx)
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

func (c *configCertificateConfigurer) Certificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	ctx := info.Context()
	cert, err := c.cache.Load(ctx)
	if err != nil {
		return nil, err
	}

	return cert.Get(info)
}

func (c *configCertificateConfigurer) slowLoad(ctx context.Context) (credential.Certifier, error) {
	certs, err := c.Enables(ctx)
	if err != nil {
		return nil, fmt.Errorf("服务端查询证书错误: %w", err)
	}

	pool := credential.NewPool()
	for _, cert := range certs {
		if err = c.parseCertificate(pool, cert.PublicKey, cert.PrivateKey); err != nil {
			return nil, err
		}
	}

	return pool, nil
}

func (c *configCertificateConfigurer) forget(enabled bool, err error) (bool, error) {
	if err == nil && enabled {
		_, _ = c.cache.Forget()
	}

	return enabled, err
}

func (c *configCertificateConfigurer) parseCertificate(pool credential.Certifier, cert, key string) error {
	certPEMBlock, keyPEMBlock := []byte(cert), []byte(key)
	pair, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return fmt.Errorf("证书私钥验证错误: %w", err)
	}
	block, _ := pem.Decode(certPEMBlock)
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("证书格式错误: %w", err)
	}
	for _, name := range x509Cert.DNSNames {
		pool.Put(name, &pair)
	}
	for _, ip := range x509Cert.IPAddresses {
		pool.Put(ip.String(), &pair)
	}

	return nil
}
