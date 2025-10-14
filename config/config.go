package config

import (
	"log/slog"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Active   string   `json:"active"`
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

type Database struct {
	URI string `json:"uri" validate:"required"`
}

type Server struct {
	Addr              string            `json:"addr"`
	ReadTimeout       Duration          `json:"read_timeout"        validate:"gte=0"`
	ReadHeaderTimeout Duration          `json:"read_header_timeout" validate:"gte=0"`
	WriteTimeout      Duration          `json:"write_timeout"       validate:"gte=0"`
	IdleTimeout       Duration          `json:"idle_timeout"        validate:"gte=0"`
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

type Duration time.Duration

func (d *Duration) UnmarshalText(raw []byte) error {
	du, err := time.ParseDuration(string(raw))
	if err != nil {
		return err
	}
	*d = Duration(du)

	return nil
}

func (d Duration) MarshalText() ([]byte, error) {
	du := time.Duration(d)
	s := du.String()

	return []byte(s), nil
}

func (d Duration) String() string {
	du := time.Duration(d)
	return du.String()
}
