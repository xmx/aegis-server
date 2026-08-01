package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BaseDB struct {
	db  *mongo.Database
	log *slog.Logger
}

func NewBaseDB(db *mongo.Database, log *slog.Logger) *BaseDB {
	return &BaseDB{
		db:  db,
		log: log,
	}
}

func (b *BaseDB) Database() *mongo.Database {
	return b.db
}

func (b *BaseDB) CreateIndex(ctx context.Context) error {
	return nil
}
