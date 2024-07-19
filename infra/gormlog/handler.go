package gormlog

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
)

type slogHandler struct {
	h slog.Handler
}

func (g *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return g.h.Enabled(ctx, level)
}

func (g *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	pcs := [10]uintptr{}
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

func (g *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return g.h.WithAttrs(attrs)
}

func (g *slogHandler) WithGroup(name string) slog.Handler {
	return g.h.WithGroup(name)
}
