package model

import "time"

// ConfigCertificate 服务端证书。
type ConfigCertificate struct {
	ID           int64     `json:"id,string"    gorm:"column:id;primaryKey;autoIncrement"`
	Enabled      bool      `json:"enabled"      gorm:"column:enabled"`
	CommonName   string    `json:"common_name"  gorm:"column:common_name;not null;index:idx_common_name"`
	PublicKey    string    `json:"public_key"   gorm:"column:public_key"`
	PrivateKey   string    `json:"private_key"  gorm:"column:private_key"`
	Organization []string  `json:"organization" gorm:"column:organization;type:json"`
	Country      []string  `json:"country"      gorm:"column:country;type:json"`
	Province     []string  `json:"province"     gorm:"column:province;type:json"`
	Names        []string  `json:"names"        gorm:"column:names;type:json"`
	DNSNames     []string  `json:"dns_names"    gorm:"column:dns_names;type:json"`
	IPAddresses  []string  `json:"ip_addresses" gorm:"column:ip_addresses;type:json"`
	Version      int       `json:"version"      gorm:"column:version"`
	NotBefore    time.Time `json:"not_before"   gorm:"column:not_before"`
	NotAfter     time.Time `json:"not_after"    gorm:"column:not_after"`
	UpdatedAt    time.Time `json:"updated_at"   gorm:"column:updated_at;not null;default:now(3)"`
	CreatedAt    time.Time `json:"created_at"   gorm:"column:created_at;not null;default:now(3)"`
}

func (ConfigCertificate) TableName() string {
	return "config_certificate"
}
