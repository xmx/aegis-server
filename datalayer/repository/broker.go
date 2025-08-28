package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Broker interface {
	Repository[bson.ObjectID, model.Broker, []*model.Broker]
}

func NewBroker(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Broker {
	coll := db.Collection("broker", opts...)
	repo := NewRepository[bson.ObjectID, model.Broker, []*model.Broker](coll)

	return &brokerRepo{
		Repository: repo,
	}
}

type brokerRepo struct {
	Repository[bson.ObjectID, model.Broker, []*model.Broker]
}

func (r *brokerRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "secret", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}
