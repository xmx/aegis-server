package request

import (
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-server/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

type InstallSetup struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

type Database struct {
	URI string `json:"uri" validate:"required"`
}

type Server struct {
	Addr              string            `json:"addr"                validate:"required"`
	ReadTimeout       model.Duration    `json:"read_timeout"        validate:"gte=0"`
	ReadHeaderTimeout model.Duration    `json:"read_header_timeout" validate:"gte=0"`
	WriteTimeout      model.Duration    `json:"write_timeout"       validate:"gte=0"`
	IdleTimeout       model.Duration    `json:"idle_timeout"        validate:"gte=0"`
	MaxHeaderBytes    int               `json:"max_header_bytes"    validate:"gte=0"`
	Static            map[string]string `json:"static"`
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

func (s InstallSetup) Config() *config.Config {
	return &config.Config{
		Server: config.Server{
			Addr:              s.Server.Addr,
			ReadTimeout:       s.Server.ReadTimeout,
			ReadHeaderTimeout: s.Server.ReadHeaderTimeout,
			WriteTimeout:      s.Server.WriteTimeout,
			IdleTimeout:       s.Server.IdleTimeout,
			MaxHeaderBytes:    s.Server.MaxHeaderBytes,
			Static:            s.Server.Static,
		},
		Database: config.Database{
			URI: s.Database.URI,
		},
		Logger: config.Logger{
			Level:   s.Logger.Level,
			Console: s.Logger.Console,
			Logger: &lumberjack.Logger{
				Filename:   config.LogFilename,
				MaxSize:    s.Logger.MaxSize,
				MaxAge:     s.Logger.MaxAge,
				MaxBackups: s.Logger.MaxBackups,
				LocalTime:  s.Logger.LocalTime,
				Compress:   s.Logger.Compress,
			},
		},
	}
}
