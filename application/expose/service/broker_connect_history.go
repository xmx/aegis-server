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

func NewBrokerConnectHistory(repo repository.All, log *slog.Logger) *BrokerConnectHistory {
	return &BrokerConnectHistory{
		repo: repo,
		log:  log,
	}
}

type BrokerConnectHistory struct {
	repo repository.All
	log  *slog.Logger
}

func (ach *BrokerConnectHistory) Page(ctx context.Context, req *request.PageKeywords) (*repository.Pages[model.BrokerConnectHistory, model.BrokerConnectHistories], error) {
	repo := ach.repo.BrokerConnectHistory()
	order := bson.D{
		{"tunnel_stat.connected_at", -1},
		{"tunnel_stat.disconnected_at", -1},
		{"_id", 1},
	}
	opt := options.Find().SetSort(order)

	return repo.FindPagination(ctx, bson.D{}, req.Page, req.Size, opt)
}
