package request

type ConfigCertificateCreate struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Enabled    bool   `json:"enabled"`
}

type ConfigCertificateUpdate struct {
	Int64ID
	ConfigCertificateCreate
}
