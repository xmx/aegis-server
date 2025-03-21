package logger

import (
	"fmt"
	"log/slog"

	"github.com/xmx/ship"
)

func NewShip(h slog.Handler, skip int) ship.Logger {
	sh := Skip(h, skip)
	log := slog.New(sh)

	return &shipLog{log: log}
}

type shipLog struct {
	log *slog.Logger
}

func (s *shipLog) Debug(format string, args ...any) {
	s.logf(slog.LevelDebug, format, args...)
}

func (s *shipLog) Info(format string, args ...any) {
	s.logf(slog.LevelInfo, format, args...)
}

func (s *shipLog) Warn(format string, args ...any) {
	s.logf(slog.LevelWarn, format, args...)
}

func (s *shipLog) Error(format string, args ...any) {
	s.logf(slog.LevelError, format, args...)
}

func (s *shipLog) logf(level slog.Level, format string, args ...any) {
	if !s.log.Enabled(nil, level) {
		return
	}

	size := len(args)
	if size == 0 {
		s.log.Log(nil, level, format)
		return
	}

	var not bool
	attrs := make([]slog.Attr, 0, size)
	for _, arg := range args {
		if attr, ok := arg.(slog.Attr); ok {
			attrs = append(attrs, attr)
		} else {
			not = true
			break
		}
	}
	if not {
		msg := fmt.Sprintf(format, args...)
		s.log.Log(nil, level, msg)
	} else {
		s.log.LogAttrs(nil, level, format, attrs...)
	}
}
