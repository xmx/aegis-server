package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User interface {
	Collection[model.User]
}

type user struct {
	Collection[model.User]
}

func NewUser(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) User {
	coll := NewCollection[model.User](db, opts...)

	return &user{
		Collection: coll,
	}
}

func (c *user) CreateIndex(ctx context.Context, opts ...options.Lister[options.CreateIndexesOptions]) ([]string, error) {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "provider", Value: 1}, {Key: "puid", Value: 1}}, Options: options.Index().SetUnique(true)},
	}

	return c.Indexes().CreateMany(ctx, indexes, opts...)
}
