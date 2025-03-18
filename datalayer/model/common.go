package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Duration struct {
	Nanosecond int64  `json:"nanosecond" bson:"nanosecond"` // 纳秒
	Formatted  string `json:"formatted"  bson:"formatted"`  // 格式化后的时间
}

type Operator struct {
	ID   bson.ObjectID `json:"id"   bson:"id"`   // 用户 ID
	Name string        `json:"name" bson:"name"` // 用户名
}
