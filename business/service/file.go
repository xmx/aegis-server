package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/condition"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type File interface {
	Cond() *response.Cond
	Page(ctx context.Context, req *request.PageKeywordOrder) ([]*model.GridFile, error)
}

func NewFile(qry *query.Query, log *slog.Logger) File {
	mod := new(model.GridFile)
	ctx := context.Background()
	tbl := qry.GridFile
	db := tbl.WithContext(ctx).UnderlyingDB()
	cond, _ := condition.ParseModel(db, mod, nil)

	return &fileService{
		qry:  qry,
		log:  log,
		cond: cond,
	}
}

type fileService struct {
	qry  *query.Query
	log  *slog.Logger
	cond *condition.Cond
}

func (svc *fileService) Cond() *response.Cond {
	return response.ReadCond(svc.cond)
}

func (svc *fileService) Page(ctx context.Context, req *request.PageKeywordOrder) ([]*model.GridFile, error) {
	tbl := svc.qry.GridFile
	dao := tbl.WithContext(ctx)
	orders := req.Order.Orders()

	return dao.Scopes(svc.cond.Order(orders)).Find()
}
