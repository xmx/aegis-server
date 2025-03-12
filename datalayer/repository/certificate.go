package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Certificate interface {
	Repository[model.Certificate]
}

func NewCertificate(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Certificate {
	base := newBaseRepository[model.Certificate](db, "certificate", opts...)
	return &certificateRepo{baseRepository: base}
}

type certificateRepo struct {
	*baseRepository[model.Certificate]
}

func (repo *certificateRepo) CreateIndex(ctx context.Context) error {
	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "certificate_sha256", Value: 1}}, // 按照 age 升序创建索引
		Options: options.Index().SetUnique(true),               // 非唯一索引
	}
	_, err := repo.coll.Indexes().CreateOne(ctx, idx)

	return err
}
