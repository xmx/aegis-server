package model

import (
	"context"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VictoriaMetrics struct {
	ID        bson.ObjectID `bson:"_id,omitempty"        json:"-"`
	Name      string        `bson:"name,omitempty"       json:"name"`
	Enabled   bool          `bson:"enabled"              json:"enabled"` // 同时最多只能有一个启用
	Method    string        `bson:"method"               json:"method"`
	Address   string        `bson:"address"              json:"address"`
	Header    MapHeader     `bson:"header"               json:"header"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
}

func (VictoriaMetrics) CollectionInfo() (string, string) {
	return "victoria_metrics", "Victoria Metrics 服务器配置信息"
}

type Dialer interface {
	DialContext(context.Context, string, string) (net.Conn, error)
}
