package repository

import (
	"context"
	"errors"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SystemConfig interface {
	Collection[model.SystemConfig]
	Get(ctx context.Context, opts ...options.Lister[options.FindOneOptions]) (*model.SystemConfig, error)
	GetFallback(ctx context.Context, opts ...options.Lister[options.FindOneOptions]) (*model.SystemConfig, error)
}

type systemConfig struct {
	Collection[model.SystemConfig]
}

func NewSystemConfig(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) SystemConfig {
	coll := NewCollection[model.SystemConfig](db, opts...)

	return &systemConfig{
		Collection: coll,
	}
}

func (c *systemConfig) Get(ctx context.Context, opts ...options.Lister[options.FindOneOptions]) (*model.SystemConfig, error) {
	return c.FindOne(ctx, bson.D{}, opts...)
}

func (c *systemConfig) GetFallback(ctx context.Context, opts ...options.Lister[options.FindOneOptions]) (*model.SystemConfig, error) {
	dat, err := c.Get(ctx, opts...)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if dat == nil {
		dat = new(model.SystemConfig)
	}

	return dat, nil
}
