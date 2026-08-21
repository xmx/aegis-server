package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"        json:"id"`
	Enabled   bool          `bson:"enabled"              json:"enabled"`
	Provider  string        `bson:"provider"             json:"provider"`
	PUID      string        `bson:"puid"                 json:"puid"`
	Login     string        `bson:"login"                json:"login"`
	Name      string        `bson:"name"                 json:"name,omitzero"`
	AvatarURL string        `bson:"avatar_url"           json:"avatar_url,omitzero"`
	Company   string        `bson:"company"              json:"company,omitzero"`
	Email     string        `bson:"email"                json:"email,omitzero"`
	Location  string        `bson:"location"             json:"location,omitzero"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
}

func (User) CollectionInfo() (string, string) { return "user", "系统用户表" }
