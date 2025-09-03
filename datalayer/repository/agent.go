package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Agent interface {
	Repository[bson.ObjectID, model.Agent, []*model.Agent]
}

func NewAgent(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Agent {
	coll := db.Collection("agent", opts...)
	repo := NewRepository[bson.ObjectID, model.Agent, []*model.Agent](coll)

	return &agentRepo{
		Repository: repo,
	}
}

type agentRepo struct {
	Repository[bson.ObjectID, model.Agent, []*model.Agent]
}

func (r *agentRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "machine_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}
