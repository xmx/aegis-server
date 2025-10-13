package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewBroker(repo repository.All, hub linkhub.Huber, log *slog.Logger) *Broker {
	return &Broker{
		repo: repo,
		hub:  hub,
		log:  log,
	}
}

type Broker struct {
	repo repository.All
	hub  linkhub.Huber
	log  *slog.Logger
}

func (b *Broker) Reset(ctx context.Context) error {
	filter := bson.M{"status": true}
	update := bson.M{"$set": bson.M{"status": false}}

	repo := b.repo.Broker()
	_, err := repo.UpdateMany(ctx, filter, update)

	return err
}

func (b *Broker) Create(ctx context.Context, name string) error {
	now := time.Now()
	buf := make([]byte, 50)
	_, _ = rand.Read(buf)
	buf[0] = 0xb0 // 这样 hex 字符串第一个字母是 b，对应 broker
	secret := hex.EncodeToString(buf)

	dat := &model.Broker{
		Name:      name,
		Secret:    secret,
		UpdatedAt: now,
		CreatedAt: now,
	}

	repo := b.repo.Broker()
	_, err := repo.InsertOne(ctx, dat)

	return err
}

func (b *Broker) List(ctx context.Context) (model.Brokers, error) {
	repo := b.repo.Broker()
	return repo.Find(ctx, bson.M{})
}

func (b *Broker) Kickout(id bson.ObjectID) error {
	peer := b.hub.GetByID(id)
	if peer == nil {
		return nil
	}
	mux := peer.Muxer()
	_ = mux.Close()

	return nil
}
