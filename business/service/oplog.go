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

type Oplog interface {
	Cond() *response.Cond
	Page(ctx context.Context, req *request.PageCondition) (*pagination.Result[*model.Oplog], error)
	Detail(ctx context.Context, id int64) (*model.Oplog, error)
	Delete(ctx context.Context, req *request.CondWhereInputs) error
	Write(ctx context.Context, oplog *model.Oplog) error
}

func NewOplog(qry *query.Query, log *slog.Logger) Oplog {
	mod := new(model.Oplog)
	ctx := context.Background()
	tbl := qry.Oplog
	db := tbl.WithContext(ctx).UnderlyingDB()
	ignores := []field.Expr{tbl.Body, tbl.Query, tbl.Header}
	opt := &condition.ParserOptions{IgnoreOrder: ignores, IgnoreWhere: ignores}
	cond, _ := condition.ParseModel(db, mod, opt)

	return &oplogService{
		qry:  qry,
		log:  log,
		cond: cond,
	}
}

type oplogService struct {
	qry  *query.Query
	log  *slog.Logger
	cond *condition.Cond
}

func (svc *oplogService) Cond() *response.Cond {
	return response.ReadCond(svc.cond)
}

func (svc *oplogService) Page(ctx context.Context, req *request.PageCondition) (*pagination.Result[*model.Oplog], error) {
	tbl := svc.qry.Oplog
	scope := svc.cond.Scope(req.AllInputs())
	dao := tbl.WithContext(ctx).Scopes(scope)
	cnt, err := dao.Count()
	if err != nil {
		return nil, err
	}
	pager := pagination.NewPager[*model.Oplog](req.PageSize())
	if cnt == 0 {
		empty := pager.Empty()
		return empty, nil
	}

	omits := []field.Expr{tbl.Body, tbl.Query, tbl.Header}
	dats, err := dao.Omit(omits...).Scopes(pager.Scope(cnt)).Find()
	if err != nil {
		return nil, err
	}
	ret := pager.Result(dats)

	return ret, nil
}

func (svc *oplogService) Detail(ctx context.Context, id int64) (*model.Oplog, error) {
	tbl := svc.qry.Oplog
	return tbl.WithContext(ctx).
		Where(tbl.ID.Eq(id)).
		First()
}

func (svc *oplogService) Delete(ctx context.Context, req *request.CondWhereInputs) error {
	wheres := svc.cond.CompileWheres(req.Inputs())
	if len(wheres) != 0 { // 禁止全表删除
		return nil
	}

	tbl := svc.qry.Oplog
	_, err := tbl.WithContext(ctx).
		Where(wheres...).
		Delete()

	return err
}

func (svc *oplogService) Write(ctx context.Context, oplog *model.Oplog) error {
	if oplog == nil {
		return nil
	}

	dao := svc.qry.Oplog.WithContext(ctx)

	return dao.Create(oplog)
}
