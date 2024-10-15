package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/condition"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/pagination"
	"github.com/xmx/aegis-server/datalayer/query"
	"gorm.io/gen/field"
)

type File interface {
	Cond() *response.Cond
	Page(ctx context.Context, req *request.PageKeywordCond) (*pagination.Result[*model.GridFile], error)
}

func NewFile(qry *query.Query, log *slog.Logger) File {
	mod := new(model.GridFile)
	ctx := context.Background()
	tbl := qry.GridFile
	db := tbl.WithContext(ctx).UnderlyingDB()
	opt := &condition.ParserOptions{
		IgnoreOrder: []field.Expr{tbl.Burst, tbl.SHA1, tbl.SHA256},
		IgnoreWhere: []field.Expr{tbl.Burst, tbl.SHA1, tbl.SHA256},
	}
	cond, _ := condition.ParseModel(db, mod, opt)

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

func (svc *fileService) Page(ctx context.Context, req *request.PageKeywordCond) (*pagination.Result[*model.GridFile], error) {
	tbl := svc.qry.GridFile
	scope := svc.cond.Scope(req.Inputs())
	dao := tbl.WithContext(ctx).Scopes(scope)
	cnt, err := dao.Count()
	if err != nil {
		return nil, err
	}
	pager := pagination.NewPager[*model.GridFile](req.PageSize())
	if cnt == 0 {
		empty := pager.Empty()
		return empty, nil
	}

	dats, err := dao.Scopes(pager.Scope(cnt)).Find()
	if err != nil {
		return nil, err
	}
	ret := pager.Result(dats)

	return ret, nil
}
