package model

import (
	"io"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"gopkg.in/natefinch/lumberjack.v2"
)

type SystemConfig struct {
	ID        bson.ObjectID `bson:"_id,omitempty"        json:"-"`
	Server    SCServer      `bson:"server,omitempty"     json:"server,omitzero"`
	Logger    SCLogger      `bson:"logger,omitempty"     json:"logger,omitzero"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
}

func (SystemConfig) CollectionInfo() (string, string) { return "system_config", "系统配置参数" }

type SCServer struct {
	Addr   string            `bson:"addr,omitempty"   json:"addr,omitzero"`                     // 监听地址
	Static map[string]string `bson:"static,omitempty" json:"static,omitzero" validate:"lte=50"` // 静态资源
}

type SCLogger struct {
	Level      string `bson:"level,omitempty"      json:"level,omitzero"      validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
	Console    string `bson:"console,omitempty"    json:"console,omitzero"    validate:"omitempty,oneof=stdout stderr"`
	Filename   string `bson:"filename,omitempty"   json:"filename,omitzero"`
	MaxSize    int    `bson:"maxsize,omitempty"    json:"maxsize,omitzero"    validate:"gte=0"`
	MaxAge     int    `bson:"maxage,omitempty"     json:"maxage,omitzero"     validate:"gte=0"`
	MaxBackups int    `bson:"maxbackups,omitempty" json:"maxbackups,omitzero" validate:"gte=0"`
	LocalTime  bool   `bson:"localtime"            json:"localtime"`
	Compress   bool   `bson:"compress"             json:"compress"`
}

func (l SCLogger) ConsoleFD() io.Writer {
	switch l.Console {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return nil
	}
}

func (l SCLogger) Lumber() *lumberjack.Logger {
	if l.Filename == "" {
		return nil
	}

	return &lumberjack.Logger{
		Filename:   l.Filename,
		MaxSize:    l.MaxSize,
		MaxAge:     l.MaxAge,
		MaxBackups: l.MaxBackups,
		LocalTime:  l.LocalTime,
		Compress:   l.Compress,
	}
}

func (l SCLogger) LevelVar() *slog.LevelVar {
	lvl := new(slog.LevelVar)
	_ = lvl.UnmarshalText([]byte(l.Level))

	return lvl
}
