package service

import (
	"io"
	"log/slog"

	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/library/dynwriter"
)

func NewLog(lvl *slog.LevelVar, writer dynwriter.Writer, log *slog.Logger) *Log {
	return &Log{
		lvl:    lvl,
		writer: writer,
		log:    log,
	}
}

type Log struct {
	lvl    *slog.LevelVar
	writer dynwriter.Writer
	log    *slog.Logger
}

func (l *Log) RegisterGlobal(vm jsvm.Runtime) error {
	fns := map[string]any{
		"level":    l.Level,
		"setLevel": l.SetLevel,
	}
	return vm.Runtime().Set("log", fns)
}

func (l *Log) Level() slog.Level {
	return l.lvl.Level()
}

func (l *Log) SetLevel(lvl string) error {
	return l.lvl.UnmarshalText([]byte(lvl))
}

func (l *Log) Attach(w io.Writer) bool {
	return l.writer.Attach(w)
}

func (l *Log) Detach(w io.Writer) bool {
	return l.writer.Detach(w)
}
