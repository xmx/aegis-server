package profile

import "gopkg.in/natefinch/lumberjack.v2"

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

func (c Logger) Lumber() *lumberjack.Logger {
	filename := c.Filename
	if filename == "" {
		return nil
	}

	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackup,
		LocalTime:  c.Localtime,
		Compress:   c.Compress,
	}
}
