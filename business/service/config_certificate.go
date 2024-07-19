package service

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log/slog"
	"sync"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/library/credential"
)

type ConfigCertificateService interface {
	Find(ctx context.Context, id int64) (*model.ConfigCertificate, error)
	List(ctx context.Context) ([]*model.ConfigCertificate, error)
	Create(ctx context.Context, req *request.ConfigCertificateCreate) error
	Update(ctx context.Context, req *request.ConfigCertificateUpdate) error
	Delete(ctx context.Context, id int64) error

	// Refresh 刷新证书。
	Refresh(ctx context.Context) error
}

func ConfigCertificate(pool credential.Certifier, qry *query.Query, log *slog.Logger) ConfigCertificateService {
	return &configCertificateService{
		pool:  pool,
		qry:   qry,
		log:   log,
		limit: 100,
	}
}

type configCertificateService struct {
	pool  credential.Certifier // 证书池。
	qry   *query.Query
	log   *slog.Logger
	limit int64      // 数据库最多可保存的证书数量。
	mutex sync.Mutex // 确保证书新增/修改/删除操作的安全性。
}

func (svc *configCertificateService) Find(ctx context.Context, id int64) (*model.ConfigCertificate, error) {
	tbl := svc.qry.ConfigCertificate
	return tbl.WithContext(ctx).Where(tbl.ID.Eq(id)).First()
}

func (svc *configCertificateService) List(ctx context.Context) ([]*model.ConfigCertificate, error) {
	return svc.qry.ConfigCertificate.WithContext(ctx).Find()
}

func (svc *configCertificateService) Create(ctx context.Context, req *request.ConfigCertificateCreate) error {
	enabled := req.Enabled
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, enabled)
	if err != nil {
		return err
	}

	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	err = svc.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		dao := tbl.WithContext(ctx)
		if cnt, exx := dao.Count(); exx != nil {
			return exx
		} else if cnt >= svc.limit {
			return fmt.Errorf("证书个数超限制 %d", svc.limit)
		}

		return dao.Create(dat)
	})
	if err != nil || !enabled {
		return err
	}
	if err = svc.Refresh(ctx); err != nil {
		svc.log.Error("新增证书后刷新证书出错", "error", err)
	}

	return err
}

func (svc *configCertificateService) Update(ctx context.Context, req *request.ConfigCertificateUpdate) error {
	enabled := req.Enabled
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, enabled)
	if err != nil {
		return err
	}

	svc.mutex.Lock()
	defer svc.mutex.Unlock()
	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)
	old, err := dao.Where(tbl.ID.Eq(req.ID)).First()
	if err != nil {
		return err
	}

	refresh := enabled || old.Enabled
	dat.ID, dat.CreatedAt = old.ID, old.CreatedAt
	err = dao.Save(dat)
	if err != nil || !refresh {
		return err
	}
	if err = svc.Refresh(ctx); err != nil {
		svc.log.Error("更新证书后刷新证书出错", "error", err)
	}

	return err
}

func (svc *configCertificateService) Delete(ctx context.Context, id int64) error {
	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)
	dat, err := dao.Where(tbl.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if _, err = dao.Delete(dat); err != nil || !dat.Enabled {
		return err
	}

	if err = svc.Refresh(ctx); err != nil {
		svc.log.Error("删除证书后刷新证书出错", "error", err)
	}

	return err
}

func (svc *configCertificateService) Refresh(ctx context.Context) error {
	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)
	dats, err := dao.Where(tbl.Enabled.Is(true)).Find()
	if err != nil {
		return err
	}

	certs := make([]tls.Certificate, 0, len(dats))
	for _, dat := range dats {
		cert, exx := tls.X509KeyPair([]byte(dat.PublicKey), []byte(dat.PrivateKey))
		if exx != nil {
			return exx
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		svc.log.Error("证书被清空，可能导致程序无法被外部访问。")
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
