package request

import "github.com/xmx/aegis-control/datalayer/model"

type BrokerUpsert struct {
	Name    string                `json:"name"    validate:"required,gte=3,lte=20"`
	Exposes model.ExposeAddresses `json:"exposes" validate:"gte=1,lte=20,unique,dive,required"`
	Config  model.BrokerConfig    `json:"config"`
}
