package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Program struct {
	ID          bson.ObjectID `bson:"_id"          json:"-"`
	PackageName string        `bson:"package_name" json:"package_name"` // 包名，全局唯一。
	Name        string        `bson:"name"         json:"name"`         // 软件名
	Version     string        `bson:"version"      json:"version"`      // 软件版本
	FileID      bson.ObjectID `bson:"file_id"      json:"file_id"`      // 文件 ID
	CreatedAt   time.Time     `bson:"created_at"   json:"created_at"`   // 安装时间
	UpdatedAt   time.Time     `bson:"updated_at"   json:"updated_at"`   // 更新时间
}
