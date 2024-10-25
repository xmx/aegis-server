// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/xmx/aegis-server/datalayer/model"
)

func newRoleMenu(db *gorm.DB, opts ...gen.DOOption) roleMenu {
	_roleMenu := roleMenu{}

	_roleMenu.roleMenuDo.UseDB(db, opts...)
	_roleMenu.roleMenuDo.UseModel(&model.RoleMenu{})

	tableName := _roleMenu.roleMenuDo.TableName()
	_roleMenu.ALL = field.NewAsterisk(tableName)
	_roleMenu.ID = field.NewInt64(tableName, "id")
	_roleMenu.RoleID = field.NewInt64(tableName, "role_id")
	_roleMenu.MenuID = field.NewInt64(tableName, "menu_id")

	_roleMenu.fillFieldMap()

	return _roleMenu
}

type roleMenu struct {
	roleMenuDo roleMenuDo

	ALL    field.Asterisk
	ID     field.Int64
	RoleID field.Int64
	MenuID field.Int64

	fieldMap map[string]field.Expr
}

func (r roleMenu) Table(newTableName string) *roleMenu {
	r.roleMenuDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r roleMenu) As(alias string) *roleMenu {
	r.roleMenuDo.DO = *(r.roleMenuDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *roleMenu) updateTableName(table string) *roleMenu {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt64(table, "id")
	r.RoleID = field.NewInt64(table, "role_id")
	r.MenuID = field.NewInt64(table, "menu_id")

	r.fillFieldMap()

	return r
}

func (r *roleMenu) WithContext(ctx context.Context) *roleMenuDo { return r.roleMenuDo.WithContext(ctx) }

func (r roleMenu) TableName() string { return r.roleMenuDo.TableName() }

func (r roleMenu) Alias() string { return r.roleMenuDo.Alias() }

func (r roleMenu) Columns(cols ...field.Expr) gen.Columns { return r.roleMenuDo.Columns(cols...) }

func (r *roleMenu) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *roleMenu) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 3)
	r.fieldMap["id"] = r.ID
	r.fieldMap["role_id"] = r.RoleID
	r.fieldMap["menu_id"] = r.MenuID
}

func (r roleMenu) clone(db *gorm.DB) roleMenu {
	r.roleMenuDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r roleMenu) replaceDB(db *gorm.DB) roleMenu {
	r.roleMenuDo.ReplaceDB(db)
	return r
}

type roleMenuDo struct{ gen.DO }

func (r roleMenuDo) Debug() *roleMenuDo {
	return r.withDO(r.DO.Debug())
}

func (r roleMenuDo) WithContext(ctx context.Context) *roleMenuDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r roleMenuDo) ReadDB() *roleMenuDo {
	return r.Clauses(dbresolver.Read)
}

func (r roleMenuDo) WriteDB() *roleMenuDo {
	return r.Clauses(dbresolver.Write)
}

func (r roleMenuDo) Session(config *gorm.Session) *roleMenuDo {
	return r.withDO(r.DO.Session(config))
}

func (r roleMenuDo) Clauses(conds ...clause.Expression) *roleMenuDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r roleMenuDo) Returning(value interface{}, columns ...string) *roleMenuDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r roleMenuDo) Not(conds ...gen.Condition) *roleMenuDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r roleMenuDo) Or(conds ...gen.Condition) *roleMenuDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r roleMenuDo) Select(conds ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r roleMenuDo) Where(conds ...gen.Condition) *roleMenuDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r roleMenuDo) Order(conds ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r roleMenuDo) Distinct(cols ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r roleMenuDo) Omit(cols ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r roleMenuDo) Join(table schema.Tabler, on ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r roleMenuDo) LeftJoin(table schema.Tabler, on ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r roleMenuDo) RightJoin(table schema.Tabler, on ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r roleMenuDo) Group(cols ...field.Expr) *roleMenuDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r roleMenuDo) Having(conds ...gen.Condition) *roleMenuDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r roleMenuDo) Limit(limit int) *roleMenuDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r roleMenuDo) Offset(offset int) *roleMenuDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r roleMenuDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *roleMenuDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r roleMenuDo) Unscoped() *roleMenuDo {
	return r.withDO(r.DO.Unscoped())
}

func (r roleMenuDo) Create(values ...*model.RoleMenu) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r roleMenuDo) CreateInBatches(values []*model.RoleMenu, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r roleMenuDo) Save(values ...*model.RoleMenu) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r roleMenuDo) First() (*model.RoleMenu, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleMenu), nil
	}
}

func (r roleMenuDo) Take() (*model.RoleMenu, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleMenu), nil
	}
}

func (r roleMenuDo) Last() (*model.RoleMenu, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleMenu), nil
	}
}

func (r roleMenuDo) Find() ([]*model.RoleMenu, error) {
	result, err := r.DO.Find()
	return result.([]*model.RoleMenu), err
}

func (r roleMenuDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RoleMenu, err error) {
	buf := make([]*model.RoleMenu, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r roleMenuDo) FindInBatches(result *[]*model.RoleMenu, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r roleMenuDo) Attrs(attrs ...field.AssignExpr) *roleMenuDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r roleMenuDo) Assign(attrs ...field.AssignExpr) *roleMenuDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r roleMenuDo) Joins(fields ...field.RelationField) *roleMenuDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r roleMenuDo) Preload(fields ...field.RelationField) *roleMenuDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r roleMenuDo) FirstOrInit() (*model.RoleMenu, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleMenu), nil
	}
}

func (r roleMenuDo) FirstOrCreate() (*model.RoleMenu, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleMenu), nil
	}
}

func (r roleMenuDo) FindByPage(offset int, limit int) (result []*model.RoleMenu, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r roleMenuDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r roleMenuDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r roleMenuDo) Delete(models ...*model.RoleMenu) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *roleMenuDo) withDO(do gen.Dao) *roleMenuDo {
	r.DO = *do.(*gen.DO)
	return r
}
