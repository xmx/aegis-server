package request

type LogLevel struct {
	Level string `json:"level" validate:"required,oneof=DEBUG INFO WARN ERROR"`
}
