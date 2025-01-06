package request

import "mime/multipart"

type ConfigCertificateCreate struct {
	PublicKey  *multipart.FileHeader `form:"public_key"  validate:"required"`
	PrivateKey *multipart.FileHeader `form:"private_key" validate:"required"`
	Enabled    bool                  `form:"enabled"`
}

type ConfigCertificateUpdate struct {
	Int64ID
	PublicKey  *multipart.FileHeader `form:"public_key"`
	PrivateKey *multipart.FileHeader `form:"private_key"`
	Enabled    bool                  `form:"enabled"`
}
