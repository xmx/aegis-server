package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/xmx/aegis-control/contract/linkhub"
	"github.com/xmx/aegis-control/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (alv *Alive) Ping(ctx context.Context, p linkhub.Peer) error {
	now := time.Now()
	id := p.ObjectID()
	filter := bson.M{"_id": id, "status": true}
	update := bson.M{"$set": bson.M{"alive_at": now}}

	repo := alv.repo.Broker()
	opt := options.FindOneAndUpdate().SetProjection(bson.M{"name": 1})
	dat, err := repo.FindOneAndUpdate(ctx, filter, update, opt)
	if err != nil {
		return err
	}
	attrs := []any{slog.Any("id", id), slog.String("name", dat.Name)}
	alv.log.Debug("broker 发来了心跳", attrs...)

	return err
}
