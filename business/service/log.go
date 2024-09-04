package service

import (
	"io"
	"log/slog"

	"github.com/xmx/aegis-server/library/ioext"
)

type Log interface {
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
