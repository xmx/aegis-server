package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/pagination"
	"github.com/xmx/aegis-server/datalayer/repository"
)

func NewOplog(repo repository.Oplog, log *slog.Logger) *Oplog {
	return &Oplog{
		repo: repo,
		log:  log,
	}
}

type Oplog struct {
	repo repository.Oplog
	log  *slog.Logger
}

func (l *Oplog) Cond() *response.Cond {
	cond := l.repo.Cond()
	return response.ReadCond(cond)
}

func (l *Oplog) Page(ctx context.Context, req *request.PageCondition) (*pagination.Result[*model.Oplog], error) {
	page := pagination.NewPager[*model.Oplog](req.PageSize())
	inputs := req.AllInputs()

	return l.repo.Page(ctx, page, inputs)
}

func (l *Oplog) Detail(ctx context.Context, id int64) (*model.Oplog, error) {
	return l.repo.FindByID(ctx, id)
}

func (l *Oplog) Delete(ctx context.Context, req *request.CondWhereInputs) error {
	inputs := req.Inputs()
	return l.repo.Delete(ctx, inputs)
}
