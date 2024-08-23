package service

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"log/slog"

	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/pscope"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/credential"
	"gorm.io/gen"
)

type ConfigCertificate interface {
	Page(ctx context.Context, req *request.PageKeyword) (*repository.Page[*model.ConfigCertificate], error)
	Find(ctx context.Context, ids []int64) ([]*model.ConfigCertificate, error)
	Create(ctx context.Context, req *request.ConfigCertificateCreate) error
	Update(ctx context.Context, req *request.ConfigCertificateUpdate) error
	Delete(ctx context.Context, id int64) error

	// Refresh 刷新证书。
	Refresh(ctx context.Context) error
}

func NewConfigCertificate(pool credential.Certifier, repo repository.ConfigCertificate, log *slog.Logger) ConfigCertificate {
	return &configCertificateService{
		pool:  pool,
		repo:  repo,
		log:   log,
		limit: 100,
	}
}

type configCertificateService struct {
	pool  credential.Certifier // 证书池。
	repo  repository.ConfigCertificate
	log   *slog.Logger
	limit int64 // 数据库最多可保存的证书数量。
}

func (svc *configCertificateService) Page(ctx context.Context, req *request.PageKeyword) (*repository.Page[*model.ConfigCertificate], error) {
	qry := svc.repo.Query()
	tbl := qry.ConfigCertificate
	cond := make([]gen.Condition, 0, 2)
	if like := req.Like(); like != "" {
		cond = append(cond, tbl.CommonName.Like(like))
	}

	return svc.repo.Page(ctx, cond, pscope.From(req.Page))
}

func (svc *configCertificateService) Find(ctx context.Context, ids []int64) ([]*model.ConfigCertificate, error) {
	return svc.repo.FindIDs(ctx, ids)
}

func (svc *configCertificateService) Create(ctx context.Context, req *request.ConfigCertificateCreate) error {
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, req.Enabled)
	if err != nil {
		return err
	}

	overflow, enabled, err := svc.repo.Create(ctx, dat, svc.limit)
	if overflow {
		return errcode.FmtTooManyCertificate.Fmt(svc.limit)
	}
	if err != nil || !enabled {
		return err
	}

	return svc.Refresh(ctx)
}

func (svc *configCertificateService) Update(ctx context.Context, req *request.ConfigCertificateUpdate) error {
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, req.Enabled)
	if err != nil {
		return err
	}

	dat.ID = req.ID
	enabled, err := svc.repo.Update(ctx, dat)
	if err != nil || !enabled {
		return err
	}

	return svc.Refresh(ctx)
}

func (svc *configCertificateService) Delete(ctx context.Context, id int64) error {
	enabled, err := svc.repo.Delete(ctx, id)
	if err != nil || !enabled {
		return err
	}

	return svc.Refresh(ctx)
}

func (svc *configCertificateService) Refresh(ctx context.Context) error {
	dats, err := svc.repo.Enables(ctx)
	if err != nil {
		svc.log.Error("查询所有开启的证书错误", slog.Any("error", err))
		return err
	}

	certs := make([]tls.Certificate, 0, len(dats))
	for _, dat := range dats {
		cert, exx := tls.X509KeyPair([]byte(dat.PublicKey), []byte(dat.PrivateKey))
		if exx != nil {
			svc.log.Error("处理证书错误", slog.Any("error", err))
			return exx
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		svc.log.Error("当前证书表未启用任何证书，程序将无法通过网络访问。")
	}
	svc.pool.Replace(certs)

	return nil
}

func (svc *configCertificateService) parseCertificate(publicKey, privateKey string, enabled bool) (*model.ConfigCertificate, error) {
	publicKeyBlock, privateKeyBlock := []byte(publicKey), []byte(privateKey)
	cert, err := tls.X509KeyPair(publicKeyBlock, privateKeyBlock)
	if err != nil {
		return nil, err
	}

	leaf := cert.Leaf
	sub := leaf.Subject
	ips := make([]string, 0, len(leaf.IPAddresses))
	for _, addr := range leaf.IPAddresses {
		ips = append(ips, addr.String())
	}

	// 计算指纹
	certSHA256, pubKeySHA256, priKeySHA256 := svc.fingerprintSHA256(cert)
	dat := &model.ConfigCertificate{
		Enabled:           enabled,
		CommonName:        sub.CommonName,
		PublicKey:         publicKey,
		PrivateKey:        privateKey,
		CertificateSHA256: certSHA256,
		PublicKeySHA256:   pubKeySHA256,
		PrivateKeySHA256:  priKeySHA256,
		Organization:      sub.Organization,
		Country:           sub.Country,
		Province:          sub.Province,
		DNSNames:          leaf.DNSNames,
		IPAddresses:       ips,
		Version:           leaf.Version,
		NotBefore:         leaf.NotBefore,
		NotAfter:          leaf.NotAfter,
	}

	return dat, nil
}

func (*configCertificateService) fingerprintSHA256(cert tls.Certificate) (certSHA256, pubKeySHA256, priKeySHA256 string) {
	leaf := cert.Leaf
	sum256 := sha256.Sum256(leaf.Raw)
	certSHA256 = hex.EncodeToString(sum256[:])

	if pkix, _ := x509.MarshalPKIXPublicKey(leaf.PublicKey); len(pkix) != 0 {
		sum := sha256.Sum256(pkix)
		pubKeySHA256 = hex.EncodeToString(sum[:])
	}

	if pkcs8, _ := x509.MarshalPKCS8PrivateKey(cert.PrivateKey); len(pkcs8) != 0 {
		sum := sha256.Sum256(pkcs8)
		priKeySHA256 = hex.EncodeToString(sum[:])
	}

	return
}
