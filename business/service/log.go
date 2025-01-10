package service

import (
	"io"
	"log/slog"

	"github.com/xmx/aegis-server/library/multiwrite"
)

func NewLog(logLevel, gormLevel *slog.LevelVar, writer multiwrite.Writer, log *slog.Logger) *Log {
	return &Log{
		logLevel:  logLevel,
		gormLevel: gormLevel,
		writer:    writer,
		log:       log,
	}
}

type Log struct {
	logLevel  *slog.LevelVar
	gormLevel *slog.LevelVar
	writer    multiwrite.Writer
	log       *slog.Logger
}

func (l *Log) Level() (slog.Level, slog.Level) {
	return l.logLevel.Level(), l.gormLevel.Level()
}

func (l *Log) SetLevel(lvl string) error {
	return l.logLevel.UnmarshalText([]byte(lvl))
}

func (l *Log) Attach(w io.Writer) bool {
	return l.writer.Attach(w)
}

func (l *Log) Detach(w io.Writer) bool {
	return l.writer.Detach(w)
}
