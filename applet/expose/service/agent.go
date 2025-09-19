package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewAgent(repo repository.All, log *slog.Logger) *Agent {
	return &Agent{
		repo: repo,
		log:  log,
	}
}

type Agent struct {
	repo repository.All
	log  *slog.Logger
}

func (a *Agent) List(ctx context.Context) ([]*model.Agent, error) {
	repo := a.repo.Agent()
	return repo.Find(ctx, bson.M{})
}
