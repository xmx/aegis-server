package repository

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OTelClient interface {
	Collection[model.OTelClient]
}

type otelClient struct {
	Collection[model.OTelClient]
}

func NewOTelClient(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) OTelClient {
	coll := NewCollection[model.OTelClient](db, opts...)

	return &otelClient{
		Collection: coll,
	}
}
