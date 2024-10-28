package service

import (
	"io"
	"log/slog"

	"github.com/xmx/aegis-server/library/ioext"
)

func NewLog(level *slog.LevelVar, writer ioext.AttachWriter, log *slog.Logger) *Log {
	return &Log{
		level:  level,
		writer: writer,
		log:    log,
	}
}

type Log struct {
	level  *slog.LevelVar
	writer ioext.AttachWriter
	log    *slog.Logger
}

func (l *Log) Level() slog.Level {
	return l.level.Level()
}

func (l *Log) SetLevel(lvl string) error {
	return l.level.UnmarshalText([]byte(lvl))
}

func (l *Log) Attach(w io.Writer) bool {
	return l.writer.Attach(w)
}

func (l *Log) Leave(w io.Writer) bool {
	return l.writer.Leave(w)
}
