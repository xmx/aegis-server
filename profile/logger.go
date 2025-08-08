package profile

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Level   string `json:"level"   validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
	Console bool   `json:"console"`
	*lumberjack.Logger
}

func (c Logger) Lumber() *lumberjack.Logger {
	if c.Logger == nil || c.Filename == "" {
		return nil
	}

	return c.Logger
}

func (c Logger) LevelVar() *slog.LevelVar {
	lv := new(slog.LevelVar)
	_ = lv.UnmarshalText([]byte(c.Level))

	return lv
}
