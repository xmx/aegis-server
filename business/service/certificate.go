package service

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"io"
	"log/slog"
	"sync"

	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/dynsql"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/library/credential"
)

func NewCertificate(pool credential.Certifier, qry *query.Query, log *slog.Logger) (*Certificate, error) {
	opt := dynsql.Options{}
	tbl, err := dynsql.Parse(qry, []any{model.Certificate{}}, opt)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		pool:  pool,
		qry:   qry,
		tbl:   tbl,
		log:   log,
		limit: 100,
	}, nil
}

type Certificate struct {
	pool  credential.Certifier // 证书池。
	log   *slog.Logger
	qry   *query.Query
	tbl   *dynsql.Table
	mutex sync.Mutex
	limit int64 // 数据库最多可保存的证书数量。
}

func (crt *Certificate) Cond() *response.Cond {
	return response.ReadCond(crt.tbl)
}

func (crt *Certificate) Page(ctx context.Context, req *request.Pages) (*response.Pages[*model.Certificate], error) {
	//tbl := crt.qry.Certificate
	//scope := crt.cond.Scope(req.AllInputs())
	//dao := tbl.WithContext(ctx).Scopes(scope)
	//cnt, err := dao.Count()
	//if err != nil {
	//	return nil, err
	//}
	//
	//pages := response.NewPages[*model.Certificate](req.PageSize())
	//if cnt == 0 {
	//	return pages.Empty(), nil
	//}
	//
	//omits := []field.Expr{tbl.PublicKey, tbl.PrivateKey}
	//dats, err := dao.Omit(omits...).Scopes(pages.FP(cnt)).Find()
	//if err != nil {
	//	return nil, err
	//}
	//
	//return pages.SetRecords(dats), nil
	return nil, nil
}

func (crt *Certificate) Find(ctx context.Context, ids []int64) ([]*model.Certificate, error) {
	if len(ids) == 0 {
		return []*model.Certificate{}, nil
	}
	tbl := crt.qry.Certificate
	dao := tbl.WithContext(ctx)
	return dao.Where(tbl.ID.In(ids...)).Find()
}

//goland:noinspection GoUnhandledErrorResult
func (crt *Certificate) Create(ctx context.Context, req *request.ConfigCertificateCreate) error {
	pubKey, priKey := req.PublicKey, req.PrivateKey
	pubKeyFile, err := pubKey.Open()
	if err != nil {
		return err
	}
	defer pubKeyFile.Close()
	priKeyFile, err := priKey.Open()
	if err != nil {
		return err
	}
	defer priKeyFile.Close()

	pubKeyBlock, err := io.ReadAll(pubKeyFile)
	if err != nil {
		return err
	}
	priKeyBlock, err := io.ReadAll(priKeyFile)
	if err != nil {
		return err
	}

	enabled := req.Enabled
	dat, err := crt.parseCertificate(pubKeyBlock, priKeyBlock, enabled)
	if err != nil {
		return err
	}

	crt.mutex.Lock()
	defer crt.mutex.Unlock()

	// 检查证书指纹，避免出现证书重复。
	tbl := crt.qry.Certificate
	dao := tbl.WithContext(ctx)
	if cnt, _ := dao.Where(tbl.CertificateSHA256.Eq(dat.CertificateSHA256)).
		Count(); cnt > 0 {
		return errcode.ErrCertificateExisted
	}

	// 检查证书是否已超过限制个数。
	if cnt, _ := dao.Count(); cnt >= crt.limit {
		return errcode.FmtTooManyCertificate.Fmt(crt.limit)
	}
	// 新增证书。
	if err = dao.Create(dat); err != nil || !enabled {
		return err
	}
	_, err = crt.Refresh(ctx)

	return err
}

func (crt *Certificate) Update(ctx context.Context, req *request.ConfigCertificateUpdate) error {
	//id, enabled := req.ID, req.Enabled
	//dat, err := svc.parseCertificate(req.PublicKey, req.PrivateKey, enabled)
	//if err != nil {
	//	return err
	//}
	//
	//tbl := svc.qry.Certificate
	//dao := tbl.WithContext(ctx)
	//
	//svc.mutex.Lock()
	//defer svc.mutex.Unlock()
	//
	//// 查询数据库中的数据。
	//mod, err := dao.Select(tbl.Enabled, tbl.CertificateSHA256).
	//	Where(tbl.ID.Eq(id)).
	//	First()
	//if err != nil {
	//	return err
	//}
	//// 指纹变了说明修改了证书
	//if mod.CertificateSHA256 != dat.CertificateSHA256 {
	//	if cnt, _ := dao.Where(tbl.CertificateSHA256.Eq(dat.CertificateSHA256)).
	//		Count(); cnt > 0 {
	//		return errcode.ErrCertificateExisted
	//	}
	//}
	//
	//enabled = enabled || mod.Enabled
	//dat.ID, dat.CreatedAt = id, mod.CreatedAt
	//if err = dao.Where(tbl.ID.Eq(id)).Save(dat); err != nil || !enabled {
	//	return err
	//}
	//_, err = svc.Refresh(ctx)
	//
	//return err
	return nil
}

func (crt *Certificate) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	crt.mutex.Lock()
	defer crt.mutex.Unlock()

	tbl := crt.qry.Certificate
	dao := tbl.WithContext(ctx)
	cnt, _ := dao.Where(tbl.ID.In(ids...), tbl.Enabled.Is(true)).Count()
	_, err := dao.Where(tbl.ID.In(ids...)).Delete()
	if err != nil || cnt == 0 {
		return err
	}
	_, err = crt.Refresh(ctx)

	return err
}

func (crt *Certificate) Detail(ctx context.Context, id int64) (*model.Certificate, error) {
	tbl := crt.qry.Certificate
	return tbl.WithContext(ctx).
		Where(tbl.ID.Eq(id)).
		First()
}

func (crt *Certificate) Refresh(ctx context.Context) (int, error) {
	tbl := crt.qry.Certificate
	dao := tbl.WithContext(ctx)
	dats, err := dao.Where(tbl.Enabled.Is(true)).Find()
	if err != nil {
		crt.log.Error("查询所有开启的证书错误", slog.Any("error", err))
		return 0, err
	}

	certs := make([]tls.Certificate, 0, len(dats))
	for _, dat := range dats {
		cert, exx := tls.X509KeyPair(dat.PublicKey, dat.PrivateKey)
		if exx != nil {
			crt.log.Error("处理证书错误", slog.Any("error", err))
			return 0, exx
		}
		certs = append(certs, cert)
	}
	num := len(certs)
	if num == 0 {
		crt.log.Error("当前证书表未启用任何证书，程序将无法通过网络访问。")
	}
	crt.pool.Replace(certs)

	return num, nil
}

func (crt *Certificate) parseCertificate(publicKey, privateKey []byte, enabled bool) (*model.Certificate, error) {
	cert, err := tls.X509KeyPair(publicKey, privateKey)
	if err != nil {
		crt.log.Warn("证书解析错误", slog.Any("error", err))
		return nil, errcode.ErrCertificateInvalid
	}

	leaf := cert.Leaf
	sub := leaf.Subject
	ips := make([]string, 0, len(leaf.IPAddresses))
	for _, addr := range leaf.IPAddresses {
		ips = append(ips, addr.String())
	}
	uris := make([]string, 0, len(leaf.URIs))
	for _, uri := range leaf.URIs {
		uris = append(uris, uri.String())
	}

	// 计算指纹
	certSHA256, pubKeySHA256, priKeySHA256 := crt.fingerprintSHA256(cert)
	dat := &model.Certificate{
		ID:                0,
		Enabled:           enabled,
		CommonName:        sub.CommonName,
		PublicKey:         publicKey,
		PrivateKey:        privateKey,
		CertificateSHA256: certSHA256,
		PublicKeySHA256:   pubKeySHA256,
		PrivateKeySHA256:  priKeySHA256,
		DNSNames:          leaf.DNSNames,
		IPAddresses:       ips,
		EmailAddresses:    leaf.EmailAddresses,
		URIs:              uris,
		Version:           leaf.Version,
		NotBefore:         leaf.NotBefore,
		NotAfter:          leaf.NotAfter,
		Issuer:            crt.parsePKIX(leaf.Issuer),
		Subject:           crt.parsePKIX(leaf.Subject),
	}

	return dat, nil
}

func (*Certificate) parsePKIX(v pkix.Name) model.PKIXName {
	return model.PKIXName{
		Country:            v.Country,
		Organization:       v.Organization,
		OrganizationalUnit: v.OrganizationalUnit,
		Locality:           v.Locality,
		Province:           v.Province,
		StreetAddress:      v.StreetAddress,
		PostalCode:         v.PostalCode,
		SerialNumber:       v.SerialNumber,
		CommonName:         v.CommonName,
	}
}

// fingerprintSHA256 计算证书和私钥的 SHA256 指纹。
func (*Certificate) fingerprintSHA256(cert tls.Certificate) (certSHA256, pubKeySHA256, priKeySHA256 string) {
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
