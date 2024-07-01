package slogx

import (
	"context"
	"log/slog"
	"runtime"
)

func Skip(h slog.Handler, n int) slog.Handler {
	if n == 0 {
		return h
	}

	return &skipHandle{
		Handler: h,
		n:       n,
	}
}

type skipHandle struct {
	slog.Handler
	n int
}

func (s *skipHandle) Handle(ctx context.Context, record slog.Record) error {
	var pcs [1]uintptr
	runtime.Callers(s.n, pcs[:])
	record.PC = pcs[0]

	return s.Handler.Handle(ctx, record)
}
