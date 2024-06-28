package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log/slog"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/datalayer/repository"
)

type ConfigCertificateService interface {
	List(ctx context.Context) ([]*model.ConfigCertificate, error)
	Create(ctx context.Context, req *request.ConfigCertificateCreate) error
	Update(ctx context.Context, req *request.ConfigCertificateUpdate) error
	Delete(ctx context.Context, id int64) error
}

func ConfigCertificate(qry *query.Query, repo repository.ConfigCertificateRepository, log *slog.Logger) ConfigCertificateService {
	return &configCertificateService{
		qry:  qry,
		log:  log,
		repo: repo,
	}
}

type configCertificateService struct {
	qry  *query.Query
	log  *slog.Logger
	repo repository.ConfigCertificateRepository
}

func (svc *configCertificateService) List(ctx context.Context) ([]*model.ConfigCertificate, error) {
	return svc.qry.ConfigCertificate.WithContext(ctx).Limit(100).Find()
}

func (svc *configCertificateService) Create(ctx context.Context, req *request.ConfigCertificateCreate) error {
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, req.Enabled)
	if err != nil {
		return err
	}
	_, err = svc.repo.Create(ctx, dat)

	return err
}

func (svc *configCertificateService) Update(ctx context.Context, req *request.ConfigCertificateUpdate) error {
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, req.Enabled)
	if err != nil {
		return err
	}
	dat.ID = req.ID
	_, err = svc.repo.Update(ctx, dat)

	return err
}

func (svc *configCertificateService) Delete(ctx context.Context, id int64) error {
	_, err := svc.repo.Delete(ctx, id)
	return err
}

func (svc *configCertificateService) parseCertificate(publicKey, privateKey string, enabled bool) (*model.ConfigCertificate, error) {
	publicKeyBlock, privateKeyBlock := []byte(publicKey), []byte(privateKey)
	block, _ := pem.Decode(publicKeyBlock)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	if _, err = tls.X509KeyPair(publicKeyBlock, privateKeyBlock); err != nil {
		return nil, err
	}

	sub := cert.Subject
	ips := make([]string, 0, len(cert.IPAddresses))
	for _, addr := range cert.IPAddresses {
		ips = append(ips, addr.String())
	}

	dat := &model.ConfigCertificate{
		Enabled:      enabled,
		CommonName:   sub.CommonName,
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		Organization: sub.Organization,
		Country:      sub.Country,
		Province:     sub.Province,
		DNSNames:     cert.DNSNames,
		IPAddresses:  ips,
		Version:      cert.Version,
		NotBefore:    cert.NotBefore,
		NotAfter:     cert.NotAfter,
	}

	return dat, nil
}
