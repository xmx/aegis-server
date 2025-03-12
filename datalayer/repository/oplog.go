package repository

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Oplog interface {
	Repository[model.Oplog]
}

func NewOplog(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Oplog {
	return newBaseRepository[model.Oplog](db, "oplog", opts...)
}
