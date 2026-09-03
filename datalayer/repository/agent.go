package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Agent interface {
	Collection[model.Agent]
}

func NewAgent(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Agent {
	coll := NewCollection[model.Agent](db, opts...)

	return &agent{
		Collection: coll,
	}
}

type agent struct {
	Collection[model.Agent]
}

func (c *agent) CreateIndex(ctx context.Context, opts ...options.Lister[options.CreateIndexesOptions]) ([]string, error) {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "machine_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: -1}}, Options: options.Index().SetUnique(true)},
	}

	return c.Indexes().CreateMany(ctx, indexes, opts...)
}
