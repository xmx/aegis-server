package request

type PlayJS struct {
	Channel string `json:"channel" validate:"required,oneof=stdin signal"`
	Message string `json:"message" validate:"required"`
}
