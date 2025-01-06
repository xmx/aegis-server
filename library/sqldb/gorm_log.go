package sqldb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

func NewLog(writer io.Writer, cfg logger.Config) (logger.Interface, *slog.LevelVar) {
	glog := &gormLogger{
		SlowThreshold:             cfg.SlowThreshold,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		ParameterizedQueries:      cfg.ParameterizedQueries,
		writer:                    writer,
	}
	if glog.SlowThreshold <= 0 {
		glog.SlowThreshold = 200 * time.Millisecond
	}

	logLevel := new(slog.LevelVar)
	level := glog.mappingLevel(cfg.LogLevel)
	logLevel.Set(level)
	logOption := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(writer, logOption)

	glog.hand = &gormLogHandler{
		handler: handler,
		level:   logLevel,
	}
	glog.log = slog.New(glog.hand)

	return glog, logLevel
}

type gormLogger struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	writer                    io.Writer
	log                       *slog.Logger
	hand                      *gormLogHandler
}

func (gl *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	lvl := gl.mappingLevel(level)
	handler := slog.NewJSONHandler(gl.writer, &slog.HandlerOptions{AddSource: true, Level: lvl})

	return &gormLogger{
		SlowThreshold:             gl.SlowThreshold,
		IgnoreRecordNotFoundError: gl.IgnoreRecordNotFoundError,
		ParameterizedQueries:      gl.ParameterizedQueries,
		writer:                    gl.writer,
		log:                       slog.New(handler),
		hand: &gormLogHandler{
			handler: handler,
			level:   gl.hand.level,
		},
	}
}

func (gl *gormLogger) Info(ctx context.Context, msg string, data ...any) {
	if gl.hand.Enabled(ctx, slog.LevelInfo) {
		str := fmt.Sprintf(msg, data...)
		gl.log.InfoContext(ctx, str)
	}
}

func (gl *gormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if gl.hand.Enabled(ctx, slog.LevelWarn) {
		str := fmt.Sprintf(msg, data...)
		gl.log.WarnContext(ctx, str)
	}
}

func (gl *gormLogger) Error(ctx context.Context, msg string, data ...any) {
	if gl.hand.Enabled(ctx, slog.LevelError) {
		str := fmt.Sprintf(msg, data...)
		gl.log.ErrorContext(ctx, str)
	}
}

func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	elapsedStr := elapsed.String()
	switch {
	case err != nil && gl.hand.Enabled(ctx, slog.LevelError) && (!errors.Is(err, logger.ErrRecordNotFound) || !gl.IgnoreRecordNotFoundError):
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
			slog.Any("error", err),
		}
		gl.printf(ctx, slog.LevelError, attrs)
	case elapsed > gl.SlowThreshold && gl.SlowThreshold != 0 && gl.hand.Enabled(ctx, slog.LevelWarn):
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
			slog.String("threshold", gl.SlowThreshold.String()),
			slog.Bool("slowed", true),
		}
		gl.printf(ctx, slog.LevelWarn, attrs)
	case gl.hand.Enabled(ctx, slog.LevelInfo):
		sql, rows := fc()
		attrs := []slog.Attr{
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsedStr),
		}
		gl.printf(ctx, slog.LevelInfo, attrs)
	}
}

func (gl *gormLogger) printf(ctx context.Context, lvl slog.Level, attrs []slog.Attr) {
	gl.log.LogAttrs(ctx, lvl, "gorm", attrs...)
}

func (gl *gormLogger) mappingLevel(lvl logger.LogLevel) slog.Level {
	switch lvl {
	case logger.Silent, logger.Error:
		return slog.LevelError
	case logger.Warn:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

type gormLogHandler struct {
	handler slog.Handler
	level   slog.Leveler
}

func (glh *gormLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return glh.level.Level() <= level
}

func (glh *gormLogHandler) Handle(ctx context.Context, record slog.Record) error {
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

	return glh.handler.Handle(ctx, record)
}

func (glh *gormLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := glh.handler.WithAttrs(attrs)
	return &gormLogHandler{
		handler: handler,
		level:   glh.level,
	}
}

func (glh *gormLogHandler) WithGroup(name string) slog.Handler {
	handler := glh.handler.WithGroup(name)
	return &gormLogHandler{
		handler: handler,
		level:   glh.level,
	}
}
