package request

import "mime/multipart"

type FSUpload struct {
	File *multipart.FileHeader `json:"file" form:"file" validate:"required"`
}
