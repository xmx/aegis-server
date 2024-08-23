package request

type TunnelConnect struct {
	UID string `json:"uid" validate:"required,lte=100"`
}
