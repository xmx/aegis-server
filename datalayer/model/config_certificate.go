package model

import "time"

// Certificate 服务端证书。
type Certificate struct {
	ID                int64     `json:"id,string"          gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Enabled           bool      `json:"enabled"            gorm:"column:enabled;comment:是否启用"`
	CommonName        string    `json:"common_name"        gorm:"column:common_name;notnull;size:255;index;comment:公用名"`
	PublicKey         []byte    `json:"public_key"         gorm:"column:public_key;type:blob;comment:证书"`
	PrivateKey        []byte    `json:"private_key"        gorm:"column:private_key;type:blob;comment:私钥"`
	CertificateSHA256 string    `json:"certificate_sha256" gorm:"column:certificate_sha256;type:char(64);comment:证书指纹"`
	PublicKeySHA256   string    `json:"public_key_sha256"  gorm:"column:public_key_sha256;type:char(64);comment:公钥指纹"`
	PrivateKeySHA256  string    `json:"private_key_sha256" gorm:"column:private_key_sha256;type:char(64);comment:私钥指纹"`
	DNSNames          []string  `json:"dns_names"          gorm:"column:dns_names;type:json;serializer:json;comment:域名"`
	IPAddresses       []string  `json:"ip_addresses"       gorm:"column:ip_addresses;type:json;serializer:json;comment:IP"`
	EmailAddresses    []string  `json:"email_addresses"    gorm:"column:email_addresses;type:json;serializer:json;comment:EmailAddresses"`
	URIs              []string  `bson:"uris"               gorm:"column:uris;type:json;serializer:json;comment:URIs"`
	Version           int       `json:"version"            gorm:"column:version;comment:证书版本"`
	NotBefore         time.Time `json:"not_before"         gorm:"column:not_before;comment:不早于"`
	NotAfter          time.Time `json:"not_after"          gorm:"column:not_after;comment:不晚于"`
	Issuer            PKIXName  `json:"issuer"             gorm:"column:issuer;type:json;serializer:json;comment:Issuer"`
	Subject           PKIXName  `json:"subject"            gorm:"column:subject;type:json;serializer:json;comment:Subject"`
	UpdatedAt         time.Time `json:"updated_at"         gorm:"column:updated_at;notnull;default:now(3);comment:更新时间"`
	CreatedAt         time.Time `json:"created_at"         gorm:"column:created_at;notnull;default:now(3);comment:创建时间"`
}

type PKIXName struct {
	Country            []string `json:"country"`
	Organization       []string `json:"organization"`
	OrganizationalUnit []string `json:"organizational_unit"`
	Locality           []string `json:"locality"`
	Province           []string `json:"province"`
	StreetAddress      []string `json:"street_address"`
	PostalCode         []string `json:"postal_code"`
	SerialNumber       string   `json:"serial_number"`
	CommonName         string   `json:"common_name"`
}

func (Certificate) TableName() string {
	return "certificate"
}
