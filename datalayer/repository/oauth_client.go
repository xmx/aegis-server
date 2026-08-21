package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OAuthClient interface {
	Collection[model.OAuthClient]
	FindByProvider(ctx context.Context, provider string, opts ...options.Lister[options.FindOneOptions]) (*model.OAuthClient, error)
}

type oauthClient struct {
	Collection[model.OAuthClient]
}

func NewOAuthClient(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) OAuthClient {
	coll := NewCollection[model.OAuthClient](db, opts...)

	return &oauthClient{
		Collection: coll,
	}
}

func (c *oauthClient) CreateIndex(ctx context.Context, opts ...options.Lister[options.CreateIndexesOptions]) ([]string, error) {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "provider", Value: 1}}, Options: options.Index().SetUnique(true)},
	}

	return c.Indexes().CreateMany(ctx, indexes, opts...)
}

func (c *oauthClient) FindByProvider(ctx context.Context, provider string, opts ...options.Lister[options.FindOneOptions]) (*model.OAuthClient, error) {
	return c.FindOne(ctx, bson.D{{Key: "provider", Value: provider}}, opts...)
}
