package service

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"log/slog"
	"sync"

	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/gormcond"
	"github.com/xmx/aegis-server/argument/pscope"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/credential"
)

type ConfigCertificate interface {
	Cond() *response.Cond
	Page(ctx context.Context, req *request.PageKeyword) (*repository.Page[*model.ConfigCertificate], error)
	Find(ctx context.Context, ids []int64) ([]*model.ConfigCertificate, error)
	Create(ctx context.Context, req *request.ConfigCertificateCreate) error
	Update(ctx context.Context, req *request.ConfigCertificateUpdate) error
	Delete(ctx context.Context, ids []int64) error

	// Refresh 刷新证书。
	Refresh(ctx context.Context) error
}

func NewConfigCertificate(pool credential.Certifier, qry *query.Query, log *slog.Logger) ConfigCertificate {
	return &configCertificateService{
		pool:  pool,
		qry:   qry,
		log:   log,
		limit: 100,
	}
}

type configCertificateService struct {
	pool  credential.Certifier // 证书池。
	log   *slog.Logger
	qry   *query.Query
	order *gormcond.Order
	mutex sync.Mutex
	limit int64 // 数据库最多可保存的证书数量。
}

func (svc *configCertificateService) Cond() *response.Cond {
	return nil
}

func (svc *configCertificateService) Page(ctx context.Context, req *request.PageKeyword) (*repository.Page[*model.ConfigCertificate], error) {
	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx).Scopes(req.LikeScope(tbl.CommonName))
	cnt, err := dao.Count()
	if err != nil {
		return nil, err
	}

	page := pscope.From(req.Page)
	if cnt == 0 {
		return repository.PageZero[*model.ConfigCertificate](page), nil
	}

	dats, err := dao.Scopes(page.Gen(cnt)).Find()
	if err != nil {
		return nil, err
	}
	ret := repository.PageRecords(page, cnt, dats)

	return ret, nil
}

func (svc *configCertificateService) Find(ctx context.Context, ids []int64) ([]*model.ConfigCertificate, error) {
	if len(ids) == 0 {
		return []*model.ConfigCertificate{}, nil
	}
	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)
	return dao.Where(tbl.ID.In(ids...)).Find()
}

func (svc *configCertificateService) Create(ctx context.Context, req *request.ConfigCertificateCreate) error {
	enabled := req.Enabled
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, enabled)
	if err != nil {
		return err
	}

	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// 检查证书指纹，避免出现证书重复。
	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)
	if cnt, _ := dao.Where(tbl.CertificateSHA256.Eq(dat.CertificateSHA256)).
		Count(); cnt > 0 {
		return errcode.ErrCertificateExisted
	}

	// 检查证书是否已超过限制个数。
	if cnt, _ := dao.Count(); cnt >= svc.limit {
		return errcode.FmtTooManyCertificate.Fmt(svc.limit)
	}
	// 新增证书。
	if err = dao.Create(dat); err != nil || !enabled {
		return err
	}

	return svc.Refresh(ctx)
}

func (svc *configCertificateService) Update(ctx context.Context, req *request.ConfigCertificateUpdate) error {
	id, enabled := req.ID, req.Enabled
	dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, enabled)
	if err != nil {
		return err
	}

	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)

	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// 查询数据库中的数据。
	mod, err := dao.Select(tbl.Enabled, tbl.CertificateSHA256).Where(tbl.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	// 指纹变了说明修改了证书
	if mod.CertificateSHA256 != dat.CertificateSHA256 {
		if cnt, _ := dao.Where(tbl.CertificateSHA256.Eq(dat.CertificateSHA256)).
			Count(); cnt > 0 {
			return errcode.ErrCertificateExisted
		}
	}

	enabled = enabled || mod.Enabled
	dat.ID, dat.CreatedAt = id, mod.CreatedAt
	if err = dao.Where(tbl.ID.Eq(id)).Save(dat); err != nil || !enabled {
		return err
	}

	return svc.Refresh(ctx)
}

func (svc *configCertificateService) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)
	cnt, _ := dao.Where(tbl.ID.In(ids...), tbl.Enabled.Is(true)).Count()
	_, err := dao.Where(tbl.ID.In(ids...)).Delete()
	if err != nil || cnt == 0 {
		return err
	}

	return svc.Refresh(ctx)
}

func (svc *configCertificateService) Refresh(ctx context.Context) error {
	tbl := svc.qry.ConfigCertificate
	dao := tbl.WithContext(ctx)
	dats, err := dao.Where(tbl.Enabled.Is(true)).Find()
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
		svc.log.Warn("证书解析错误", slog.Any("error", err))
		return nil, errcode.ErrCertificateInvalid
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
		Locality:          sub.Locality,
		DNSNames:          leaf.DNSNames,
		IPAddresses:       ips,
		Version:           leaf.Version,
		NotBefore:         leaf.NotBefore,
		NotAfter:          leaf.NotAfter,
	}
	if dat.Organization == nil {
		dat.Organization = []string{}
	}
	if dat.Country == nil {
		dat.Country = []string{}
	}
	if dat.Province == nil {
		dat.Province = []string{}
	}
	if dat.Locality == nil {
		dat.Locality = []string{}
	}
	if dat.DNSNames == nil {
		dat.DNSNames = []string{}
	}

	return dat, nil
}

// fingerprintSHA256 计算证书和私钥的 SHA256 指纹。
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
