package logger

import (
	"bytes"
	"log"
	"log/slog"
)

func NewV1(l *slog.Logger, skip ...int) *log.Logger {
	var n int
	if len(skip) > 0 {
		n = skip[0]
	}
	if n != 0 {
		l = slog.New(Skip(l.Handler(), n))
	}
	w := &v1Writer{l: l}

	return log.New(w, "", 0)
}

type v1Writer struct {
	l *slog.Logger
}

func (v *v1Writer) Write(p []byte) (int, error) {
	n := len(p)
	s := bytes.TrimRight(p, "\n")
	v.l.Info(string(s))

	return n, nil
}
