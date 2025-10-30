package request

type InstallSetup struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

type Database struct {
	URI string `json:"uri" validate:"required"`
}

type Server struct {
	Addr   string            `json:"addr"   validate:"required"`
	Static map[string]string `json:"static"`
}

type Logger struct {
	Level      string `json:"level"      validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
	Console    bool   `json:"console"`
	MaxSize    int    `json:"maxsize"    validate:"gte=0"`
	MaxAge     int    `json:"maxage"     validate:"gte=0"`
	MaxBackups int    `json:"maxbackups" validate:"gte=0"`
	LocalTime  bool   `json:"localtime"`
	Compress   bool   `json:"compress"`
}
