package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AgentConnRecord interface {
	Collection[model.AgentConnRecord]
}

func NewAgentConnRecord(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) AgentConnRecord {
	coll := NewCollection[model.AgentConnRecord](db, opts...)

	return &agentConnRecord{
		Collection: coll,
	}
}

type agentConnRecord struct {
	Collection[model.AgentConnRecord]
}

func (c *agentConnRecord) CreateIndex(ctx context.Context, opts ...options.Lister[options.CreateIndexesOptions]) ([]string, error) {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "agent_id", Value: 1}}},
		{Keys: bson.D{{Key: "disconnected_at", Value: -1}}},
	}

	return c.Indexes().CreateMany(ctx, indexes, opts...)
}
