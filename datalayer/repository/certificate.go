package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Certificate interface {
	Repository[bson.ObjectID, model.Certificate, []*model.Certificate]
}

func NewCertificate(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Certificate {
	coll := db.Collection("certificate", opts...)
	repo := NewRepository[bson.ObjectID, model.Certificate, []*model.Certificate](coll)

	return &certificateRepo{
		Repository: repo,
	}
}

type certificateRepo struct {
	Repository[bson.ObjectID, model.Certificate, []*model.Certificate]
}

func (r *certificateRepo) CreateIndex(ctx context.Context) error {
	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "certificate_sha256", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.Indexes().CreateOne(ctx, idx)

	return err
}
