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
	"github.com/xmx/aegis-server/application/errcode"
	"github.com/xmx/aegis-server/application/expose/request"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (brk *Broker) Reset(ctx context.Context) error {
	filter := bson.M{"status": true}
	update := bson.M{"$set": bson.M{"status": false}}

	repo := brk.repo.Broker()
	_, err := repo.UpdateMany(ctx, filter, update)

	return err
}

func (brk *Broker) Page(ctx context.Context, req *request.PageKeywords) (*repository.Pages[model.Broker, model.Brokers], error) {
	repo := brk.repo.Broker()

	return repo.FindPagination(ctx, bson.D{}, req.Page, req.Size)
}

func (brk *Broker) Create(ctx context.Context, req *request.BrokerCreate) error {
	now := time.Now()
	buf := make([]byte, 50)
	_, _ = rand.Read(buf)
	buf[0] = 0xb0 // 这样 hex 字符串第一个字母是 b，对应 broker
	secret := hex.EncodeToString(buf)

	dat := &model.Broker{
		Name:      req.Name,
		Exposes:   req.Exposes,
		Secret:    secret,
		UpdatedAt: now,
		CreatedAt: now,
	}

	repo := brk.repo.Broker()
	_, err := repo.InsertOne(ctx, dat)

	return err
}

func (brk *Broker) Kickout(id bson.ObjectID) error {
	peer := brk.hub.GetByID(id)
	if peer == nil {
		return nil
	}
	mux := peer.Muxer()
	_ = mux.Close()

	return nil
}

func (brk *Broker) Exposes(ctx context.Context) (model.ExposeAddresses, error) {
	opt := options.Find().SetProjection(bson.M{"exposes": 1})
	repo := brk.repo.Broker()
	brks, err := repo.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	var exposes model.ExposeAddresses
	for _, b := range brks {
		exposes = append(exposes, b.Exposes...)
	}
	if len(exposes) == 0 {
		return nil, errcode.ErrNilDocument
	}

	return exposes, nil
}
