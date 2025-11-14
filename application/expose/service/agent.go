package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/application/expose/request"
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

func (agt *Agent) Page(ctx context.Context, req *request.PageKeywords) (*repository.Pages[model.Agent, model.Agents], error) {
	fields := []string{
		"execute_stat.goos", "execute_stat.goarch", "execute_stat.hostname",
		"machine_id", "networks.name", "networks.ipv4", "networks.ipv6",
		"tunnel_stat.local_addr", "tunnel_stat.remote_addr",
	}
	filter := make(bson.M, 4)
	if arr := req.Regexps(fields); len(arr) != 0 {
		filter["$or"] = arr
	}
	repo := agt.repo.Agent()

	return repo.FindPagination(ctx, filter, req.Page, req.Size)
}
