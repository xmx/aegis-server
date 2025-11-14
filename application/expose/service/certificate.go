package service

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"iter"
	"log/slog"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/tlscert"
	"github.com/xmx/aegis-server/application/errcode"
	"github.com/xmx/aegis-server/application/expose/request"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewCertificate(repo repository.All, pool tlscert.Certifier, log *slog.Logger) *Certificate {
	return &Certificate{
		repo: repo,
		pool: pool,
		log:  log,
	}
}

type Certificate struct {
	repo repository.All
	pool tlscert.Certifier
	log  *slog.Logger
}

func (crt *Certificate) Page(ctx context.Context, req *request.PageKeywords) (*repository.Pages[model.Certificate, model.Certificates], error) {
	filter := make(bson.M, 4)
	fields := []string{"common_name", "dns_names"}
	if arr := req.Regexps(fields); len(arr) != 0 {
		filter["$or"] = arr
	}
	repo := crt.repo.Certificate()

	proj := bson.M{"public_key": 0, "private_key": 0}
	opt := options.Find().SetProjection(proj)

	return repo.FindPagination(ctx, filter, req.Page, req.Size, opt)
}

func (crt *Certificate) Get(ctx context.Context, id bson.ObjectID) (*model.Certificate, error) {
	repo := crt.repo.Certificate()
	return repo.FindByID(ctx, id)
}

func (crt *Certificate) Create(ctx context.Context, req *request.ConfigCertificateCreate) error {
	now := time.Now()
	pubKey, priKey, enabled := req.PublicKey, req.PrivateKey, req.Enabled
	mod, err := crt.Parse(pubKey, priKey)
	if err != nil {
		return err
	}
	mod.Name = req.Name
	mod.Enabled = enabled
	mod.UpdatedAt = now
	mod.CreatedAt = now

	repo := crt.repo.Certificate()
	if _, err = repo.InsertOne(ctx, mod); err != nil {
		return err
	}
	if enabled {
		crt.Reset(ctx)
	}

	return nil
}

func (crt *Certificate) Update(ctx context.Context, req *request.ConfigCertificateUpdate) error {
	id := req.OID()
	now := time.Now()
	pubKey, priKey, enabled := req.PublicKey, req.PrivateKey, req.Enabled
	mod, err := crt.Parse(pubKey, priKey)
	if err != nil {
		return err
	}
	mod.Name = req.Name
	mod.Enabled = enabled
	mod.UpdatedAt = now

	repo := crt.repo.Certificate()
	proj := bson.M{"public_key": 0, "private_key": 0}
	opt := options.FindOneAndUpdate().SetProjection(proj)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": mod}
	last, err := repo.FindOneAndUpdate(ctx, filter, update, opt)
	if err != nil {
		return err
	}
	if enabled || last.Enabled {
		crt.Reset(ctx)
	}

	return nil
}

// Delete 通过 ID 删除证书。
func (crt *Certificate) Delete(ctx context.Context, ids []bson.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}

	repo := crt.repo.Certificate()
	filter := bson.M{"_id": bson.M{"$in": ids}}
	res, err := repo.DeleteMany(ctx, filter)
	if err != nil {
		return err
	} else if res.DeletedCount == 0 {
		return errcode.ErrNilDocument
	}
	crt.Reset(ctx)

	return nil
}

// All 遍历证书，如果 ID 为空则代表查询所有证书。
func (crt *Certificate) All(ctx context.Context, ids []bson.ObjectID) iter.Seq2[*model.Certificate, error] {
	filter := make(bson.M, 4)
	if len(ids) != 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	repo := crt.repo.Certificate()

	return repo.All(ctx, filter)
}

func (crt *Certificate) Reset(ctx context.Context) {
	crt.log.WarnContext(ctx, "清除证书缓存")
	crt.pool.Reset()
}

func (crt *Certificate) Parse(publicKey, privateKey string) (*model.Certificate, error) {
	pubKey, priKey := []byte(publicKey), []byte(privateKey)
	cert, err := tls.X509KeyPair(pubKey, priKey)
	if err != nil {
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
		CommonName:         sub.CommonName,
		PublicKey:          publicKey,
		PrivateKey:         privateKey,
		CertificateSHA256:  certSHA256,
		PublicKeySHA256:    pubKeySHA256,
		PrivateKeySHA256:   priKeySHA256,
		DNSNames:           leaf.DNSNames,
		IPAddresses:        ips,
		EmailAddresses:     leaf.EmailAddresses,
		URIs:               uris,
		Version:            leaf.Version,
		NotBefore:          leaf.NotBefore,
		NotAfter:           leaf.NotAfter,
		Issuer:             crt.parsePKIX(leaf.Issuer),
		Subject:            crt.parsePKIX(leaf.Subject),
		SignatureAlgorithm: cert.Leaf.SignatureAlgorithm.String(),
	}

	return dat, nil
}

func (*Certificate) parsePKIX(v pkix.Name) model.CertificatePKIXName {
	return model.CertificatePKIXName{
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

	return certSHA256, pubKeySHA256, priKeySHA256
}
