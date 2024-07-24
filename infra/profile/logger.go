package profile

import (
	"io"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Level     string `json:"level"`
	Terminal  bool   `json:"terminal"`
	Filename  string `json:"filename"`
	MaxAge    int    `json:"max_age"`
	MaxSize   int    `json:"max_size"`
	MaxBackup int    `json:"max_backup"`
	Localtime bool   `json:"localtime"`
	Compress  bool   `json:"compress"`
}

func (c Logger) Writer() LogWriteCloser {
	lwc := c.newLogWriter()
	terminal, filename := c.Terminal, c.Filename
	if !terminal && filename == "" {
		return lwc
	}

	if terminal {
		lwc.Append(&nopWriteCloser{w: os.Stdout})
	}

	if filename != "" {
		lumber := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackup,
			LocalTime:  c.Localtime,
			Compress:   c.Compress,
		}
		lwc.Append(lumber)
	}

	return lwc
}

func (Logger) newLogWriter() LogWriteCloser {
	return &logWriteCloser{elems: make(map[io.WriteCloser]struct{}, 8)}
}
