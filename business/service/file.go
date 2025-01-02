package service

import (
	"context"
	"io"
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

func NewFile(qry *query.Query, dbfs gridfs.FS, log *slog.Logger) *File {
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
		dbfs: dbfs,
		cond: cond,
	}
}

type File struct {
	qry  *query.Query
	log  *slog.Logger
	dbfs gridfs.FS
	cond *condition.Cond
}

func (f *File) Cond() *response.Cond {
	return response.ReadCond(f.cond)
}

func (f *File) Page(ctx context.Context, req *request.PageCondition) (*pagination.Result[*model.GridFile], error) {
	tbl := f.qry.GridFile
	scope := f.cond.Scope(req.AllInputs())
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

func (f *File) Count(ctx context.Context, limit int) (response.NameCounts, error) {
	if limit <= 0 {
		limit = 1
	}

	ret := make(response.NameCounts, 0, limit)
	nameExpr, countExpr := ret.Aliases()
	nameAlias := nameExpr.ColumnName().String()
	countAlias := countExpr.ColumnName().String()

	tbl := f.qry.GridFile
	err := tbl.WithContext(ctx).
		Select(tbl.Extension.As(nameAlias), tbl.Extension.Count().As(countAlias)).
		Group(tbl.Extension).
		Order(countExpr.Desc(), tbl.Extension).
		Limit(limit).
		Scan(&ret)

	return ret, err
}

func (f *File) Save(ctx context.Context, filename string, r io.Reader) (*model.GridFile, error) {
	return f.dbfs.Save(ctx, filename, r)
}

func (f *File) Open(ctx context.Context, fileID int64) (gridfs.File, error) {
	return f.dbfs.OpenID(ctx, fileID)
}
