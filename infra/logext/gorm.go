package logext

import (
	"context"
	"errors"
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

func (l *gormLog) LogMode(lvl logger.LogLevel) logger.Interface {
	nl := *l
	nl.LogLevel = lvl
	return &nl
}

func (l *gormLog) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		l.log.InfoContext(ctx, msg, slog.Any("data", data))
	}
}

func (l *gormLog) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		l.log.WarnContext(ctx, msg, slog.Any("data", data))
	}
}

func (l *gormLog) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		l.log.ErrorContext(ctx, msg, slog.Any("data", data))
	}
}

func (l *gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedStr := elapsed.String()
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
			slog.Any("error", err),
		}
		l.printf(ctx, slog.LevelError, attrs)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
			slog.String("threshold", l.SlowThreshold.String()),
			slog.Bool("slowed", true),
		}
		l.printf(ctx, slog.LevelWarn, attrs)
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
		}
		l.printf(ctx, slog.LevelInfo, attrs)
	}
}

func (l *gormLog) printf(ctx context.Context, lvl slog.Level, attrs []slog.Attr) {
	l.log.LogAttrs(ctx, lvl, "gorm", attrs...)
}

type gormHandler struct {
	h slog.Handler
}

func (g *gormHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return g.h.Enabled(ctx, level)
}

func (g *gormHandler) Handle(ctx context.Context, record slog.Record) error {
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

	return g.h.Handle(ctx, record)
}

func (g *gormHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return g.h.WithAttrs(attrs)
}

func (g *gormHandler) WithGroup(name string) slog.Handler {
	return g.h.WithGroup(name)
}
