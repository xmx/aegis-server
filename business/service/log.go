package service

import (
	"io"
	"log/slog"

	"github.com/xmx/aegis-server/library/multiwrite"
)

func NewLog(lvl *slog.LevelVar, writer multiwrite.Writer, log *slog.Logger) *Log {
	return &Log{
		lvl:    lvl,
		writer: writer,
		log:    log,
	}
}

type Log struct {
	lvl    *slog.LevelVar
	writer multiwrite.Writer
	log    *slog.Logger
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
