package logext

import (
	"fmt"
	"log/slog"

	"github.com/xgfone/ship/v5"
)

func Ship(h slog.Handler) ship.Logger {
	sh := Skip(h, 6)
	log := slog.New(sh)

	return &shipLog{log: log}
}

type shipLog struct {
	log *slog.Logger
}

func (s *shipLog) Tracef(format string, args ...any) {
	s.logf(slog.LevelDebug, format, args...)
}

func (s *shipLog) Debugf(format string, args ...any) {
	s.logf(slog.LevelDebug, format, args...)
}

func (s *shipLog) Infof(format string, args ...any) {
	s.logf(slog.LevelInfo, format, args...)
}

func (s *shipLog) Warnf(format string, args ...any) {
	s.logf(slog.LevelWarn, format, args...)
}

func (s *shipLog) Errorf(format string, args ...any) {
	s.logf(slog.LevelError, format, args...)
}

func (s *shipLog) logf(level slog.Level, format string, args ...any) {
	if s.log.Enabled(nil, level) {
		msg := fmt.Sprintf(format, args...)
		s.log.Log(nil, level, msg)
	}
}
