package request

import "mime/multipart"

type AgentReleaseUpload struct {
	File      *multipart.FileHeader `json:"file"      form:"file"      validate:"required"`
	Changelog string                `json:"changelog" form:"changelog" validate:"lte=10000"`
}
