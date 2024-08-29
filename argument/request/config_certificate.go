package request

type ConfigCertificateCreate struct {
	PublicKey  string `json:"public_key"  validate:"required"`
	PrivateKey string `json:"private_key" validate:"required"`
	Enabled    bool   `json:"enabled"`
}

type ConfigCertificateUpdate struct {
	Int64ID
	ConfigCertificateCreate
}
