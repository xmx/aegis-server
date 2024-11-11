package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/condition"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/pagination"
	"github.com/xmx/aegis-server/datalayer/query"
	"gorm.io/gen/field"
)

type Oplog interface {
	Cond() *condition.Cond
	Create(ctx context.Context, mod *model.Oplog) error
	FindByID(ctx context.Context, id int64) (*model.Oplog, error)
	Page(ctx context.Context, page pagination.Pager[*model.Oplog], scope *condition.ScopeInput) (*pagination.Result[*model.Oplog], error)
	Delete(ctx context.Context, where *condition.WhereInputs) error
}

func NewOplog(qry *query.Query) Oplog {
	mod := new(model.Oplog)
	ctx := context.Background()
	tbl := qry.Oplog
	db := tbl.WithContext(ctx).UnderlyingDB()
	ignores := []field.Expr{tbl.Body, tbl.Query, tbl.Header}
	opt := &condition.ParserOptions{IgnoreOrder: ignores, IgnoreWhere: ignores}
	cond, _ := condition.ParseModel(db, mod, opt)

	return &oplogRepository{
		qry:  qry,
		cond: cond,
	}
}

type oplogRepository struct {
	qry  *query.Query
	cond *condition.Cond
}

func (op *oplogRepository) Cond() *condition.Cond {
	return op.cond
}

func (op *oplogRepository) Create(ctx context.Context, mod *model.Oplog) error {
	mod.ID = 0 // reset id
	tbl := op.qry.Oplog

	return tbl.WithContext(ctx).Create(mod)
}

func (op *oplogRepository) FindByID(ctx context.Context, id int64) (*model.Oplog, error) {
	tbl := op.qry.Oplog
	return tbl.WithContext(ctx).
		Where(tbl.ID.Eq(id)).
		First()
}

func (op *oplogRepository) Page(ctx context.Context, page pagination.Pager[*model.Oplog], scope *condition.ScopeInput) (*pagination.Result[*model.Oplog], error) {
	search := op.cond.Scope(scope)
	tbl := op.qry.Oplog
	dao := tbl.WithContext(ctx).Scopes(search)
	cnt, err := dao.Count()
	if err != nil {
		return nil, err
	} else if cnt <= 0 {
		empty := page.Empty()
		return empty, nil
	}

	dats, err := dao.Scopes(page.Scope(cnt)).Find()
	if err != nil {
		return nil, err
	}
	ret := page.Result(dats)

	return ret, nil
}

func (op *oplogRepository) Delete(ctx context.Context, where *condition.WhereInputs) error {
	wheres := op.cond.CompileWheres(where)
	if len(wheres) == 0 { // 禁止全表删除
		return nil
	}
	tbl := op.qry.Oplog
	_, err := tbl.WithContext(ctx).
		Where(wheres...).
		Delete()

	return err
}
