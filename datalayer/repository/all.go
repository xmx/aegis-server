package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type All interface {
	DB() *mongo.Database
	Client() *mongo.Client

	Agent() Agent
	Broker() Broker
	Certificate() Certificate

	GridFSBucket(opts ...options.Lister[options.BucketOptions]) *mongo.GridFSBucket
	CreateIndex(ctx context.Context) error
}

func NewAll(db *mongo.Database) All {
	return allRepo{
		db:          db,
		agent:       NewAgent(db),
		broker:      NewBroker(db),
		certificate: NewCertificate(db),
	}
}

type allRepo struct {
	db          *mongo.Database
	agent       Agent
	broker      Broker
	certificate Certificate
}

func (ar allRepo) DB() *mongo.Database   { return ar.db }
func (ar allRepo) Client() *mongo.Client { return ar.db.Client() }

func (ar allRepo) Agent() Agent             { return ar.agent }
func (ar allRepo) Broker() Broker           { return ar.broker }
func (ar allRepo) Certificate() Certificate { return ar.certificate }

func (ar allRepo) GridFSBucket(opts ...options.Lister[options.BucketOptions]) *mongo.GridFSBucket {
	return ar.db.GridFSBucket(opts...)
}

func (ar allRepo) CreateIndex(ctx context.Context) error {
	fields := []any{
		ar.agent,
		ar.broker,
		ar.certificate,
	}
	for _, f := range fields {
		idx, ok := f.(CreateIndexer)
		if !ok {
			continue
		}
		if err := idx.CreateIndex(ctx); err != nil {
			return err
		}
	}

	return nil
}

type CreateIndexer interface {
	CreateIndex(context.Context) error
}
