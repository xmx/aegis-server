package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Operator struct {
	ID   bson.ObjectID `json:"id"   bson:"id"`   // 用户 ID
	Name string        `json:"name" bson:"name"` // 用户名
}
