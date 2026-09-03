package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel"
)

type OTelClient struct {
	db  *repository.BaseDB
	log *slog.Logger
}

func NewOTelClient(db *repository.BaseDB, log *slog.Logger) *OTelClient {
	return &OTelClient{
		db:  db,
		log: log,
	}
}

func (oc *OTelClient) List(ctx context.Context) ([]*model.OTelClient, error) {
	coll := oc.db.OTelClient()
	return coll.Find(ctx, bson.D{})
}

func (oc *OTelClient) Create(ctx context.Context, req *request.OTelClientCreate) error {
	now := time.Now()
	dat := &model.OTelClient{
		Enabled:   req.Enabled,
		Endpoint:  req.Endpoint,
		Protocol:  req.Protocol,
		Insecure:  req.Insecure,
		Header:    req.Header,
		CreatedAt: now,
		UpdatedAt: now,
	}
	coll := oc.db.OTelClient()
	_, err := coll.InsertOne(ctx, dat)

	return err
}

func (oc *OTelClient) Update(ctx context.Context, req *request.OTelClientUpdate) error {
	id, now := req.ID.UnsafeID(), time.Now()
	dat := &model.OTelClient{
		Enabled:   req.Enabled,
		Endpoint:  req.Endpoint,
		Protocol:  req.Protocol,
		Insecure:  req.Insecure,
		Header:    req.Header,
		UpdatedAt: now,
	}

	update := bson.M{"$set": dat}

	coll := oc.db.OTelClient()
	_, err := coll.UpdateByID(ctx, id, update)

	return err
}

func (oc *OTelClient) Delete(ctx context.Context, id bson.ObjectID) error {
	tracer := otel.Tracer("aegis-server.service")
	_, span := tracer.Start(ctx, "otel_client.delete")
	defer span.End()

	coll := oc.db.OTelClient()
	dat, err := coll.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	_ = dat

	return nil
}
