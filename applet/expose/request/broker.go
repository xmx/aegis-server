package request

import "github.com/xmx/aegis-control/datalayer/model"

type BrokerCreate struct {
	Name   string             `json:"name"   validate:"required,gte=3,lte=20"`
	Config model.BrokerConfig `json:"config"`
}
