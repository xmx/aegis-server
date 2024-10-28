package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/condition"
	"github.com/xmx/aegis-server/datalayer/gridfs"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/pagination"
	"github.com/xmx/aegis-server/datalayer/query"
	"gorm.io/gen/field"
)

func NewFile(qry *query.Query, log *slog.Logger) *File {
	mod := new(model.GridFile)
	ctx := context.Background()
	tbl := qry.GridFile
	db := tbl.WithContext(ctx).UnderlyingDB()
	ignores := []field.Expr{tbl.Burst, tbl.MD5, tbl.SHA1, tbl.SHA256}
	opt := &condition.ParserOptions{IgnoreOrder: ignores, IgnoreWhere: ignores}
	cond, _ := condition.ParseModel(db, mod, opt)

	return &File{
		qry:  qry,
		log:  log,
		cond: cond,
	}
}

type File struct {
	qry  *query.Query
	log  *slog.Logger
	dbfs gridfs.File
	cond *condition.Cond
}

func (svc *File) Cond() *response.Cond {
	return response.ReadCond(svc.cond)
}

func (svc *File) Page(ctx context.Context, req *request.PageCondition) (*pagination.Result[*model.GridFile], error) {
	tbl := svc.qry.GridFile
	scope := svc.cond.Scope(req.AllInputs())
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

func (svc *File) Count(ctx context.Context, limit int) (response.NameCounts, error) {
	if limit <= 0 {
		limit = 1
	}

	ret := make(response.NameCounts, 0, limit)
	nameAlias, countAlias, countField := ret.Aliases()

	tbl := svc.qry.GridFile
	err := tbl.WithContext(ctx).
		Select(tbl.Extension.As(nameAlias), tbl.Extension.Count().As(countAlias)).
		Group(tbl.Extension).
		Order(countField.Desc(), tbl.Extension).
		Limit(limit).
		Scan(&ret)
	return ret, err
}
