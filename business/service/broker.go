package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewBroker(repo repository.All, log *slog.Logger) *Broker {
	return &Broker{
		repo: repo,
		log:  log,
	}
}

type Broker struct {
	repo repository.All
	log  *slog.Logger
}

func (b *Broker) Reset(ctx context.Context) error {
	filter := bson.M{"status": true}
	update := bson.M{"$set": bson.M{"status": false}}

	repo := b.repo.Broker()
	_, err := repo.UpdateMany(ctx, filter, update)

	return err
}
