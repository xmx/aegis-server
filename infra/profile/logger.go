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
		lwc.Append(&nopCloseWriter{w: os.Stdout})
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

//func (c Logger) Option() (*slog.HandlerOptions, io.WriteCloser) {
//	terminal, filename := c.Terminal, c.Filename
//	if !terminal && filename == "" {
//		opt := &slog.HandlerOptions{Level: slog.LevelError}
//		return opt, new(loggerDiscord)
//	}
//
//	var closer io.Closer
//	writers := make([]io.Writer, 0, 2)
//	if terminal {
//		writers = append(writers, os.Stdout)
//	}
//	if filename != "" {
//		lumber := &lumberjack.Logger{
//			Filename:   filename,
//			MaxSize:    c.MaxSize,
//			MaxAge:     c.MaxAge,
//			MaxBackups: c.MaxBackup,
//			LocalTime:  c.Localtime,
//			Compress:   c.Compress,
//		}
//		closer = lumber
//		writers = append(writers, lumber)
//	}
//
//	var w io.Writer
//	if len(writers) == 1 {
//		w = writers[0]
//	} else {
//		w = io.MultiWriter(writers...)
//	}
//	wrt := &loggerWriter{w: w, c: closer}
//
//	level := slog.LevelInfo
//	_ = level.UnmarshalText([]byte(c.Level))
//	opt := &slog.HandlerOptions{Level: level, AddSource: true}
//
//	return opt, wrt
//}

//
//type loggerDiscord struct{}
//
//func (l *loggerDiscord) Write(p []byte) (int, error) { return len(p), nil }
//func (l *loggerDiscord) Close() error                { return nil }
//
//type loggerWriter struct {
//	w io.Writer
//	c io.Closer
//}
//
//func (l loggerWriter) Write(p []byte) (int, error) {
//	return l.w.Write(p)
//}
//
//func (l loggerWriter) Close() error {
//	if c := l.c; c != nil {
//		return c.Close()
//	}
//
//	return nil
//}
