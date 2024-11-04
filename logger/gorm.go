package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

func Gorm(h slog.Handler, cfg logger.Config) logger.Interface {
	lh := &gormHandler{h: h}
	log := slog.New(lh)

	threshold := cfg.SlowThreshold
	if threshold <= 0 {
		threshold = 200 * time.Millisecond
	}
	level := cfg.LogLevel
	if level == 0 {
		level = logger.Warn
	}

	return &gormLog{
		SlowThreshold:             threshold,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		ParameterizedQueries:      cfg.ParameterizedQueries,
		LogLevel:                  level,
		log:                       log,
	}
}

type gormLog struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  logger.LogLevel
	log                       *slog.Logger
}

func (gl *gormLog) LogMode(lvl logger.LogLevel) logger.Interface {
	nl := *gl
	nl.LogLevel = lvl
	return &nl
}

func (gl *gormLog) Info(ctx context.Context, msg string, data ...any) {
	if gl.LogLevel >= logger.Info {
		str := fmt.Sprintf(msg, data...)
		gl.log.InfoContext(ctx, str)
	}
}

func (gl *gormLog) Warn(ctx context.Context, msg string, data ...any) {
	if gl.LogLevel >= logger.Warn {
		str := fmt.Sprintf(msg, data...)
		gl.log.WarnContext(ctx, str)
	}
}

func (gl *gormLog) Error(ctx context.Context, msg string, data ...any) {
	if gl.LogLevel >= logger.Error {
		str := fmt.Sprintf(msg, data...)
		gl.log.ErrorContext(ctx, str)
	}
}

func (gl *gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if gl.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedStr := elapsed.String()
	switch {
	case err != nil && gl.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !gl.IgnoreRecordNotFoundError):
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
			slog.Any("error", err),
		}
		gl.printf(ctx, slog.LevelError, attrs)
	case elapsed > gl.SlowThreshold && gl.SlowThreshold != 0 && gl.LogLevel >= logger.Warn:
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
			slog.String("threshold", gl.SlowThreshold.String()),
			slog.Bool("slowed", true),
		}
		gl.printf(ctx, slog.LevelWarn, attrs)
	case gl.LogLevel == logger.Info:
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
		}
		gl.printf(ctx, slog.LevelInfo, attrs)
	}
}

func (gl *gormLog) printf(ctx context.Context, lvl slog.Level, attrs []slog.Attr) {
	gl.log.LogAttrs(ctx, lvl, "gorm", attrs...)
}

type gormHandler struct {
	h slog.Handler
}

func (gh *gormHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return gh.h.Enabled(ctx, level)
}

func (gh *gormHandler) Handle(ctx context.Context, record slog.Record) error {
	// https://github.com/go-gorm/gorm/blob/v1.25.12/utils/utils.go#L33-L49
	pcs := [13]uintptr{}
	size := runtime.Callers(6, pcs[:])
	frames := runtime.CallersFrames(pcs[:size])
	for i := 0; i < size; i++ {
		frame, _ := frames.Next()
		file := frame.File
		if (!strings.HasPrefix(file, "gorm.io/") || strings.HasSuffix(file, "_test.go")) &&
			!strings.HasSuffix(file, ".gen.go") {
			record.PC = pcs[i]
			break
		}
	}

	return gh.h.Handle(ctx, record)
}

func (gh *gormHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return gh.h.WithAttrs(attrs)
}

func (gh *gormHandler) WithGroup(name string) slog.Handler {
	return gh.h.WithGroup(name)
}
