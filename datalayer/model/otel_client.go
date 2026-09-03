package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OTelClient struct {
	ID        bson.ObjectID `json:"id"         bson:"_id,omitempty"`
	Enabled   bool          `json:"enabled"    bson:"enabled"`
	Endpoint  string        `json:"endpoint"   bson:"endpoint"`
	Protocol  string        `json:"protocol"   bson:"protocol"` // grpc http
	Insecure  bool          `json:"insecure"   bson:"insecure"`
	Header    MapHeader     `json:"header"     bson:"header"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at,omitempty"`
}

func (OTelClient) CollectionInfo() (string, string) { return "otel_client", "OTEL 配置" }
