package logger

import (
	"context"
	"fmt"
	"log/slog"
)

type Format struct {
	log *slog.Logger
}

// NewFormat format 式的日志输出。
// ship: 6
// pyroscope: 5-6
func NewFormat(h slog.Handler, skip int) *Format {
	sh := Skip(h, skip)

	return &Format{
		log: slog.New(sh),
	}
}

func (f *Format) Tracef(format string, args ...any) {
	f.logf(slog.LevelDebug, format, args)
}

func (f *Format) Debugf(format string, args ...any) {
	f.logf(slog.LevelDebug, format, args)
}

func (f *Format) Infof(format string, args ...any) {
	f.logf(slog.LevelInfo, format, args)
}

func (f *Format) Warnf(format string, args ...any) {
	f.logf(slog.LevelWarn, format, args)
}

func (f *Format) Errorf(format string, args ...any) {
	f.logf(slog.LevelError, format, args)
}

func (f *Format) logf(level slog.Level, format string, args []any) {
	ctx := context.Background()
	if !f.log.Enabled(ctx, level) {
		return
	}

	if len(args) != 0 {
		format = fmt.Sprintf(format, args...)
	}

	f.log.Log(ctx, level, format)
}
