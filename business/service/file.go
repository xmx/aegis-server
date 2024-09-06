package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/argument/bizdata"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type File interface {
	Page(ctx context.Context, req *request.PageKeywordOrder) ([]*model.GridFile, error)
}

func NewFile(qry *query.Query, log *slog.Logger) File {
	tbl := qry.GridFile
	orders := new(bizdata.SearchOrders).
		Add(tbl.Filename, "文件名").
		Add(tbl.Length, "文件大小").
		Add(tbl.CreatedAt, "上传时间")

	return &fileService{
		qry:    qry,
		log:    log,
		orders: orders,
	}
}

type fileService struct {
	qry    *query.Query
	log    *slog.Logger
	orders *bizdata.SearchOrders
}

func (svc *fileService) Page(ctx context.Context, req *request.PageKeywordOrder) ([]*model.GridFile, error) {
	tbl := svc.qry.GridFile
	dao := tbl.WithContext(ctx)

	return dao.Scopes(svc.orders.Scope(req.Order)).Find()
}
