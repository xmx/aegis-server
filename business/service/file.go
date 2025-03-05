package service

import (
	"context"
	"io"
	"log/slog"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/datalayer/dynsql"
	"github.com/xmx/aegis-server/datalayer/gridfs"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"gorm.io/gen/field"
)

func NewFile(qry *query.Query, dbfs gridfs.FS, log *slog.Logger) (*File, error) {
	fl := &File{qry: qry, log: log, dbfs: dbfs}
	opt := dynsql.Options{Where: fl.whereFunc(), Order: fl.orderFunc()}
	mods := []any{model.GridFile{}}

	tbl, err := dynsql.Parse(qry, mods, opt)
	if err != nil {
		return nil, err
	}
	fl.tbl = tbl

	return fl, nil
}

type File struct {
	qry  *query.Query
	log  *slog.Logger
	dbfs gridfs.FS
	tbl  *dynsql.Table
}

func (f *File) Register(root *jsvm.Object) error {
	root.Sub("service").Set("file", f)
	return nil
}

func (f *File) Cond() *response.Cond {
	return response.ReadCond(f.tbl)
}

func (f *File) Page(ctx context.Context, req *request.Pages) (*response.Pages[*model.GridFile], error) {
	//tbl := f.qry.GridFile
	//scope := f.cond.Scope(req.AllInputs())
	//dao := tbl.WithContext(ctx).Scopes(scope)
	//cnt, err := dao.Count()
	//if err != nil {
	//	return nil, err
	//}
	//pager := pagination.NewPager[*model.GridFile](req.PageSize())
	//if cnt == 0 {
	//	empty := pager.Empty()
	//	return empty, nil
	//}
	//
	//dats, err := dao.Scopes(pager.Scope(cnt)).Find()
	//if err != nil {
	//	return nil, err
	//}
	//ret := pager.Result(dats)
	//
	//return ret, nil
	return nil, nil
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

func (f *File) whereFunc() func(*dynsql.Where) *dynsql.Where {
	file := f.qry.GridFile
	ignores := []field.Expr{file.MD5, file.SHA1, file.SHA256}
	return func(where *dynsql.Where) *dynsql.Where {
		for _, ignore := range ignores {
			if where.Equals(ignore) {
				return nil
			}
		}
		return where
	}
}

func (f *File) orderFunc() func(*dynsql.Order) *dynsql.Order {
	file := f.qry.GridFile
	ignores := []field.Expr{file.MD5, file.SHA1, file.SHA256, file.MediaType, file.Burst}
	return func(order *dynsql.Order) *dynsql.Order {
		for _, ignore := range ignores {
			if order.Equals(ignore) {
				return nil
			}
		}
		return order
	}
}
