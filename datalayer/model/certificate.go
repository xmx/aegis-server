package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Certificate 服务端证书。
type Certificate struct {
	ID                 bson.ObjectID       `json:"id,omitzero"           bson:"_id,omitempty"`
	Name               string              `json:"name,omitzero"         bson:"name"`
	Enabled            bool                `json:"enabled,omitzero"      bson:"enabled"`
	CommonName         string              `json:"common_name"           bson:"common_name"`
	PublicKey          string              `json:"public_key,omitempty"  bson:"public_key"`
	PrivateKey         string              `json:"private_key,omitempty" bson:"private_key"`
	CertificateSHA256  string              `json:"certificate_sha256"    bson:"certificate_sha256"`
	PublicKeySHA256    string              `json:"public_key_sha256"     bson:"public_key_sha256"`
	PrivateKeySHA256   string              `json:"private_key_sha256"    bson:"private_key_sha256"`
	DNSNames           []string            `json:"dns_names"             bson:"dns_names"`
	IPAddresses        []string            `json:"ip_addresses"          bson:"ip_addresses"`
	EmailAddresses     []string            `json:"email_addresses"       bson:"email_addresses"`
	URIs               []string            `json:"uris"                  bson:"uris"`
	Version            int                 `json:"version"               bson:"version"`
	NotBefore          time.Time           `json:"not_before"            bson:"not_before"`
	NotAfter           time.Time           `json:"not_after"             bson:"not_after"`
	Issuer             CertificatePKIXName `json:"issuer"                bson:"issuer"`
	Subject            CertificatePKIXName `json:"subject"               bson:"subject"`
	SignatureAlgorithm string              `json:"signature_algorithm"   bson:"signature_algorithm"`
	UpdatedAt          time.Time           `json:"updated_at,omitzero"   bson:"updated_at,omitempty"`
	CreatedAt          time.Time           `json:"created_at,omitzero"   bson:"created_at,omitempty"`
}

type CertificatePKIXName struct {
	Country            []string `json:"country"             bson:"country"`
	Organization       []string `json:"organization"        bson:"organization"`
	OrganizationalUnit []string `json:"organizational_unit" bson:"organizational_unit"`
	Locality           []string `json:"locality"            bson:"locality"`
	Province           []string `json:"province"            bson:"province"`
	StreetAddress      []string `json:"street_address"      bson:"street_address"`
	PostalCode         []string `json:"postal_code"         bson:"postal_code"`
	SerialNumber       string   `json:"serial_number"       bson:"serial_number"`
	CommonName         string   `json:"common_name"         bson:"common_name"`
}
