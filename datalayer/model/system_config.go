package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SystemConfig struct {
	ID        bson.ObjectID `bson:"_id,omitempty"        json:"-"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
}

func (SystemConfig) CollectionInfo() (string, string) { return "system_config", "系统配置参数" }

type SCLogger struct {
}

type SCServer struct {
	Addr   string            `bson:"addr"   json:"addr"`                     // 监听地址
	Static map[string]string `bson:"static" json:"static" validate:"lte=50"` // 静态资源
}
