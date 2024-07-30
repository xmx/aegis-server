package request

type PlayJS struct {
	Script string `json:"script" validate:"required"`
}
