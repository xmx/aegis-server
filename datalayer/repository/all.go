package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type All interface {
	DB() *mongo.Database
	Client() *mongo.Client
	Certificate() Certificate
	Oplog() Oplog
	CreateIndex(ctx context.Context) error
}

func NewAll(db *mongo.Database) All {
	return allRepo{
		db:          db,
		certificate: NewCertificate(db),
		oplog:       NewOplog(db),
	}
}

type allRepo struct {
	db          *mongo.Database
	certificate Certificate
	oplog       Oplog
}

func (ar allRepo) DB() *mongo.Database      { return ar.db }
func (ar allRepo) Client() *mongo.Client    { return ar.db.Client() }
func (ar allRepo) Certificate() Certificate { return ar.certificate }
func (ar allRepo) Oplog() Oplog             { return ar.oplog }

func (ar allRepo) CreateIndex(ctx context.Context) error {
	indexes := []IndexCreator{
		ar.certificate, ar.oplog,
	}
	for _, idx := range indexes {
		if err := idx.CreateIndex(ctx); err != nil {
			return err
		}
	}

	return nil
}
