package model

import "time"

// ConfigCertificate 服务端证书。
type ConfigCertificate struct {
	ID                int64     `json:"id,string"          gorm:"column:id;primaryKey;autoIncrement;comment:表ID"`
	Enabled           bool      `json:"enabled"            gorm:"column:enabled;comment:是否启用"`
	CommonName        string    `json:"common_name"        gorm:"column:common_name;not null;type:varchar(255);index:idx_common_name;comment:公用名"`
	PublicKey         string    `json:"public_key"         gorm:"column:public_key;comment:证书"`
	PrivateKey        string    `json:"private_key"        gorm:"column:private_key;comment:私钥"`
	CertificateSHA256 string    `json:"certificate_sha256" gorm:"column:certificate_sha256;type:char(64);comment:证书指纹"`
	PublicKeySHA256   string    `json:"public_key_sha256"  gorm:"column:public_key_sha256;type:char(64);comment:公钥指纹"`
	PrivateKeySHA256  string    `json:"private_key_sha256" gorm:"column:private_key_sha256;type:char(64);comment:私钥指纹"`
	Organization      []string  `json:"organization"       gorm:"column:organization;type:json;serializer:json;comment:组织"`
	Country           []string  `json:"country"            gorm:"column:country;type:json;serializer:json;comment:国家"`
	Province          []string  `json:"province"           gorm:"column:province;type:json;serializer:json"`
	Locality          []string  `json:"locality"           gorm:"column:locality;type:json;serializer:json"`
	DNSNames          []string  `json:"dns_names"          gorm:"column:dns_names;type:json;serializer:json"`
	IPAddresses       []string  `json:"ip_addresses"       gorm:"column:ip_addresses;type:json;serializer:json"`
	Version           int       `json:"version"            gorm:"column:version"`
	NotBefore         time.Time `json:"not_before"         gorm:"column:not_before"`
	NotAfter          time.Time `json:"not_after"          gorm:"column:not_after"`
	UpdatedAt         time.Time `json:"updated_at"         gorm:"column:updated_at;not null;default:now(3)"`
	CreatedAt         time.Time `json:"created_at"         gorm:"column:created_at;not null;default:now(3)"`
}

func (ConfigCertificate) TableName() string {
	return "config_certificate"
}
