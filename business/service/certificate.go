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
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/credential"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewCertificate(repo repository.Certificate, pool credential.Certifier, log *slog.Logger) (*Certificate, error) {
	return &Certificate{
		repo: repo,
		pool: pool,
		log:  log,
	}, nil
}

type Certificate struct {
	repo repository.Certificate
	pool credential.Certifier // 证书池。
	log  *slog.Logger
}

func (crt *Certificate) Page(ctx context.Context, req *request.PageKeywords) (*repository.Pages[model.Certificate], error) {
	filter := make(bson.M, 4)
	if arr := req.Regexps("common_name", "dns_names"); len(arr) != 0 {
		filter["$or"] = arr
	}

	projection := bson.M{"public_key": 0, "private_key": 0}
	opt := options.Find().SetProjection(projection)

	return crt.repo.FindPage(ctx, filter, req.Page, req.Size, opt)
}

func (crt *Certificate) Find(ctx context.Context, ids []bson.ObjectID) ([]*model.Certificate, error) {
	if len(ids) == 0 {
		return []*model.Certificate{}, nil
	}

	return crt.repo.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
}

//goland:noinspection GoUnhandledErrorResult
func (crt *Certificate) Create(ctx context.Context, req *request.ConfigCertificateCreate) error {
	now := time.Now()
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

	// 检查证书指纹，避免出现证书重复。
	filter := bson.M{"certificate_sha256": dat.CertificateSHA256}
	if cnt, exx := crt.repo.CountDocuments(ctx, filter); exx != nil {
		return exx
	} else if cnt > 0 {
		return errcode.ErrCertificateExisted
	}
	dat.CreatedAt, dat.UpdatedAt = now, now
	_, err = crt.repo.InsertOne(ctx, dat)

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

func (crt *Certificate) Delete(ctx context.Context, ids []bson.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	if _, err := crt.repo.DeleteMany(ctx, filter); err != nil {
		return err
	}
	_, err := crt.Refresh(ctx)

	return err
}

func (crt *Certificate) Detail(ctx context.Context, id bson.ObjectID) (*model.Certificate, error) {
	return crt.repo.FindByID(ctx, id)
}

func (crt *Certificate) Refresh(ctx context.Context) (int, error) {
	crts, err := crt.repo.Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return 0, err
	} else if len(crts) == 0 {
		return 0, nil
	}

	certs := make([]tls.Certificate, 0, len(crts))
	for _, dat := range crts {
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

	if pki, _ := x509.MarshalPKIXPublicKey(leaf.PublicKey); len(pki) != 0 {
		sum := sha256.Sum256(pki)
		pubKeySHA256 = hex.EncodeToString(sum[:])
	}

	if pkcs8, _ := x509.MarshalPKCS8PrivateKey(cert.PrivateKey); len(pkcs8) != 0 {
		sum := sha256.Sum256(pkcs8)
		priKeySHA256 = hex.EncodeToString(sum[:])
	}

	return
}
