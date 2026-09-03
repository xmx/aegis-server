package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Agent struct {
	db  *repository.BaseDB
	log *slog.Logger
}

func NewAgent(db *repository.BaseDB, log *slog.Logger) *Agent {
	return &Agent{
		db:  db,
		log: log,
	}
}

func (agt *Agent) Page(ctx context.Context, req *request.Pages) (*repository.Pages[model.Agent], error) {
	coll := agt.db.Agent()
	order := bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: -1}}
	opts := options.Find().SetSort(order)

	return coll.Page(ctx, bson.D{}, req.Page, req.Size, opts)
}

func (agt *Agent) Get(ctx context.Context, id bson.ObjectID) (*model.Agent, error) {
	coll := agt.db.Agent()
	return coll.FindByID(ctx, id)
}

func (agt *Agent) Records(ctx context.Context, req *request.IDPages) (*repository.Pages[model.AgentConnRecord], error) {
	filter := bson.D{{Key: "agent_id", Value: req.ID.UnsafeID()}}

	coll := agt.db.AgentConnRecord()
	order := bson.D{{Key: "disconnected_at", Value: -1}}
	opts := options.Find().SetSort(order)

	return coll.Page(ctx, filter, req.Page, req.Size, opts)
}
