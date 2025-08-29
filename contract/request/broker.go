package request

type BrokerCreate struct {
	Name string `json:"name" validate:"required,gte=3,lte=20"`
}
