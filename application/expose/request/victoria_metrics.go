package request

import "github.com/xmx/aegis-control/datalayer/model"

type VictoriaMetricsUpsert struct {
	Name    string           `json:"name"    validate:"required,gte=2,lte=20"`
	Method  string           `json:"method"  validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Address string           `json:"address" validate:"http_url,lte=255"`
	Header  model.HTTPHeader `json:"header"  validate:"lte=20"`
	Enabled bool             `json:"enabled"`
}
