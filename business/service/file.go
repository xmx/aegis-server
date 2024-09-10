package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/argument/gormcond"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type File interface {
	Cond() *response.Cond
	Page(ctx context.Context, req *request.PageKeywordOrder) ([]*model.GridFile, error)
}

func NewFile(qry *query.Query, log *slog.Logger) File {
	tbl := qry.GridFile
	order := gormcond.NewOrder().
		Add(tbl.Filename, "文件名").
		Add(tbl.Length, "文件大小").
		Add(tbl.CreatedAt, "上传时间").
		Add(tbl.ID, "ID")

	return &fileService{
		qry:   qry,
		log:   log,
		order: order,
	}
}

type fileService struct {
	qry   *query.Query
	log   *slog.Logger
	order *gormcond.Order
}

func (svc *fileService) Cond() *response.Cond {
	return &response.Cond{Orders: svc.order.Columns()}
}

func (svc *fileService) Page(ctx context.Context, req *request.PageKeywordOrder) ([]*model.GridFile, error) {
	tbl := svc.qry.GridFile
	dao := tbl.WithContext(ctx)
	orders := req.Order.Orders()

	return dao.Scopes(svc.order.Scope(orders)).Find()
}
