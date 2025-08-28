package bservice

import (
	"context"
	"log/slog"
	"time"

	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewAlive(repo repository.All, log *slog.Logger) *Alive {
	return &Alive{
		repo: repo,
		log:  log,
	}
}

type Alive struct {
	repo repository.All
	log  *slog.Logger
}

func (alv *Alive) Ping(ctx context.Context, id bson.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id, "status": true}
	update := bson.M{"$set": bson.M{"alive_at": now}}

	repo := alv.repo.Broker()
	_, err := repo.UpdateOne(ctx, filter, update)

	return err
}
