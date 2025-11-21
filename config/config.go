package config

import (
	"log/slog"

	"github.com/xmx/aegis-control/datalayer/model"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// Filename 默认配置文件位置。
	Filename          = "resources/config/application.jsonc"
	LogFilename       = "resources/log/application.jsonl"
	EnvKeyInitialAddr = "AEGIS_INIT_ADDR"
)

type Config struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
	Victoria Victoria `json:"victoria"` // 临时
}

type Database struct {
	URI string `json:"uri" validate:"required"`
}

type Server struct {
	Addr              string            `json:"addr"`
	ReadTimeout       model.Duration    `json:"read_timeout"        validate:"gte=0"`
	ReadHeaderTimeout model.Duration    `json:"read_header_timeout" validate:"gte=0"`
	WriteTimeout      model.Duration    `json:"write_timeout"       validate:"gte=0"`
	IdleTimeout       model.Duration    `json:"idle_timeout"        validate:"gte=0"`
	MaxHeaderBytes    int               `json:"max_header_bytes"    validate:"gte=0"`
	Static            map[string]string `json:"static"              validate:"lte=255"`
}

type Logger struct {
	Level   string `json:"level"   validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
	Console bool   `json:"console"`
	*lumberjack.Logger
}

func (l Logger) LevelVar() *slog.LevelVar {
	lvl := new(slog.LevelVar) // default: INFO
	_ = lvl.UnmarshalText([]byte(l.Level))

	return lvl
}

func (l Logger) Lumber() *lumberjack.Logger {
	if lg := l.Logger; lg != nil && lg.Filename != "" {
		return lg
	}

	return nil
}

type Victoria struct {
	Addr   string   `json:"addr"`
	Header []string `json:"header"`
}
