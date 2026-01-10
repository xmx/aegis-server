package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/application/expose/request"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewAgentConnectHistory(repo repository.All, log *slog.Logger) *AgentConnectHistory {
	return &AgentConnectHistory{
		repo: repo,
		log:  log,
	}
}

type AgentConnectHistory struct {
	repo repository.All
	log  *slog.Logger
}

func (ach *AgentConnectHistory) Page(ctx context.Context, req *request.PageKeywords) (*repository.Pages[model.AgentConnectHistory, model.AgentConnectHistories], error) {
	repo := ach.repo.AgentConnectHistory()
	order := bson.D{
		{"tunnel_stat.connected_at", -1},
		{"tunnel_stat.disconnected_at", -1},
		{"_id", 1},
	}
	opt := options.Find().SetSort(order)

	return repo.FindPagination(ctx, bson.D{}, req.Page, req.Size, opt)
}
