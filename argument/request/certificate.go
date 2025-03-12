package request

import "mime/multipart"

type ConfigCertificateCreate struct {
	PublicKey  *multipart.FileHeader `form:"public_key"  validate:"required"`
	PrivateKey *multipart.FileHeader `form:"private_key" validate:"required"`
	Enabled    bool                  `form:"enabled"`
}

type ConfigCertificateUpdate struct {
	ObjectID
	PublicKey  *multipart.FileHeader `json:"public_key"  form:"public_key"  validate:"required"`
	PrivateKey *multipart.FileHeader `json:"private_key" form:"private_key" validate:"required"`
	Enabled    bool                  `json:"enabled"     form:"enabled"`
}
