package request

import "mime/multipart"

type PlayJS struct {
	Script string `json:"script" validate:"required"`
	Args   any    `json:"args"`
}

type PlayUpload struct {
	File *multipart.FileHeader `json:"file" form:"file" validate:"required"`
}
