package service

import (
	"io"
	"log/slog"

	"github.com/xmx/aegis-server/library/multiwrite"
)

func NewLog(level *slog.LevelVar, writer multiwrite.Writer, log *slog.Logger) *Log {
	return &Log{
		level:  level,
		writer: writer,
		log:    log,
	}
}

type Log struct {
	level  *slog.LevelVar
	writer multiwrite.Writer
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

func (l *Log) Detach(w io.Writer) bool {
	return l.writer.Detach(w)
}
