package request

type ConfigCertificateCreate struct {
	Name       string `json:"name"        validate:"required,lte=20"`
	PublicKey  string `json:"public_key"  validate:"required"`
	PrivateKey string `json:"private_key" validate:"required"`
	Enabled    bool   `json:"enabled"`
}

type ConfigCertificateUpdate struct {
	ObjectID
	Name       string `json:"name"        validate:"required,lte=20"`
	PublicKey  string `json:"public_key"  validate:"required"`
	PrivateKey string `json:"private_key" validate:"required"`
	Enabled    bool   `json:"enabled"`
}

type ConfigCertificateParse struct {
	PublicKey  string `json:"public_key"  validate:"required"`
	PrivateKey string `json:"private_key" validate:"required"`
}
