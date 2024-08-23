package service

import (
	"io"
	"log/slog"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/library/ioext"
)

type Log interface {
	jsvm.Loader
	Level() slog.Level
	SetLevel(lvl string) error
	Attach(w io.Writer) bool
	Leave(w io.Writer) bool
}

func NewLog(level *slog.LevelVar, writer ioext.AttachWriter, log *slog.Logger) Log {
	return &logService{
		level:  level,
		writer: writer,
		log:    log,
	}
}

type logService struct {
	level  *slog.LevelVar
	writer ioext.AttachWriter
	log    *slog.Logger
}

func (svc *logService) Level() slog.Level {
	return svc.level.Level()
}

func (svc *logService) SetLevel(lvl string) error {
	return svc.level.UnmarshalText([]byte(lvl))
}

func (svc *logService) Attach(w io.Writer) bool {
	return svc.writer.Attach(w)
}

func (svc *logService) Leave(w io.Writer) bool {
	return svc.writer.Leave(w)
}

func (svc *logService) Global(*goja.Runtime) error {
	return nil
}

func (svc *logService) Require() (string, require.ModuleLoader) {
	return "service/log", svc.require
}

func (svc *logService) require(_ *goja.Runtime, obj *goja.Object) {
	fields := map[string]any{
		"level":    svc.Level,
		"setLevel": svc.SetLevel,
	}
	for k, v := range fields {
		_ = obj.Set(k, v)
	}
}
