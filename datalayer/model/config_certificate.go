package model

import "time"

// ConfigCertificate 服务端证书。
type ConfigCertificate struct {
	ID                int64     `json:"id,string"          gorm:"column:id;primaryKey;autoIncrement"`
	Enabled           bool      `json:"enabled"            gorm:"column:enabled"` // 是否启用证书。
	CommonName        string    `json:"common_name"        gorm:"column:common_name;not null;index:idx_common_name"`
	PublicKey         string    `json:"public_key"         gorm:"column:public_key"`
	PrivateKey        string    `json:"private_key"        gorm:"column:private_key"`
	CertificateSHA256 string    `json:"certificate_sha256" gorm:"column:certificate_sha256;char(64)"` // 证书 SHA-256 指纹
	PublicKeySHA256   string    `json:"public_key_sha256"  gorm:"column:public_key_sha256;char(64)"`  // 公钥 SHA-256 指纹
	PrivateKeySHA256  string    `json:"private_key_sha256" gorm:"column:private_key_sha256;char(64)"` // 私钥 SHA-256 指纹
	Organization      []string  `json:"organization"       gorm:"column:organization;type:json;serializer:json"`
	Country           []string  `json:"country"            gorm:"column:country;type:json;serializer:json"`
	Province          []string  `json:"province"           gorm:"column:province;type:json;serializer:json"`
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
