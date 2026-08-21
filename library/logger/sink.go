package logger

import (
	"context"
	"log/slog"
)

func NewSink(h slog.Handler, skip int) Sink {
	han := Skip(h, skip)
	return Sink{
		log: slog.New(han),
	}
}

type Sink struct {
	log *slog.Logger
}

func (sk Sink) Info(level int, message string, keysAndValues ...any) {
	// mongo 日志输出级别定义 https://github.com/mongodb/mongo-go-driver/blob/v2.5.0/internal/logger/level.go#L23-L35
	if level == 0 { // mongo.LevelOff = 0
		return
	}

	lvl := slog.LevelDebug // mongo.LevelDebug = 2
	if level == 1 {        // mongo.LevelInfo = 1
		lvl = slog.LevelInfo
	}

	msg := sk.translate(message)
	sk.log.Log(context.Background(), lvl, msg, keysAndValues...)
}

func (sk Sink) Error(err error, message string, keysAndValues ...any) {
	kvs := []any{"err", err}
	kvs = append(kvs, keysAndValues...)
	msg := sk.translate(message)

	sk.log.Error(msg, kvs...)
}

func (Sink) translate(msg string) string {
	if msg == "Command succeeded" {
		return "命令执行成功（MongoDB 驱动）"
	} else if msg == "Command started" {
		return "命令开始执行（MongoDB 驱动）"
	} else if msg == "Command failed" {
		return "命令执行失败（MongoDB 驱动）"
	}

	return msg
}
