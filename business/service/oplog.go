package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

func NewOplog(qry *query.Query, log *slog.Logger) *Oplog {
	return &Oplog{
		qry: qry,
		log: log,
	}
}

type Oplog struct {
	qry *query.Query
	log *slog.Logger
}

func (ol *Oplog) Cond() *response.Cond {
	return nil
	// cond := ol.repo.Cond()
	// return response.ReadCond(cond)
}

func (ol *Oplog) Page(ctx context.Context, req *request.Pages) (*response.Pages[*model.Oplog], error) {
	//search := op.cond.Scope(scope)
	//tbl := op.qry.Oplog
	//dao := tbl.WithContext(ctx).Scopes(search)
	//cnt, err := dao.Count()
	//if err != nil {
	//	return nil, err
	//} else if cnt <= 0 {
	//	empty := page.Empty()
	//	return empty, nil
	//}
	//
	//dats, err := dao.Scopes(page.Scope(cnt)).Find()
	//if err != nil {
	//	return nil, err
	//}
	//ret := page.Result(dats)
	//
	//return ret, nil

	return nil, nil
}

func (ol *Oplog) Detail(ctx context.Context, id int64) (*model.Oplog, error) {
	tbl := ol.qry.Oplog
	return tbl.WithContext(ctx).Where(tbl.ID.Eq(id)).First()
}

func (ol *Oplog) Delete(ctx context.Context, req *request.Pages) error {
	// inputs := req.Inputs()
	// return ol.qry.Delete(ctx, inputs)
	return nil
}

func (ol *Oplog) Create(ctx context.Context, data *model.Oplog) error {
	data.ID = 0 // reset id
	return ol.qry.Oplog.WithContext(ctx).Create(data)
}
