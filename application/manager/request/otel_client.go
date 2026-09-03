package request

import "github.com/xmx/aegis-server/datalayer/model"

type OTelClientUpdate struct {
	ObjectID
	OTelClientCreate
}

type OTelClientCreate struct {
	Enabled  bool            `json:"enabled"`
	Endpoint string          `json:"endpoint" validate:"required"`
	Protocol string          `json:"protocol" validate:"required,oneof=http grpc"`
	Insecure bool            `json:"insecure"`
	Header   model.MapHeader `json:"header"   validate:"lte=10"`
}
