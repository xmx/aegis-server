package repository

import (
	"context"
	"sync"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"gorm.io/gen"
)

type ConfigCertificate interface {
	Repository[*model.ConfigCertificate]

	// Enables 查询已经启用的证书。
	Enables(ctx context.Context) ([]*model.ConfigCertificate, error)

	Create(ctx context.Context, cert *model.ConfigCertificate, limit int64) (overflow, enabled bool, err error)

	// Update 更新证书内容，用 ID 作为搜索条件。
	// 返回修改之前或修改之后是否启用。
	Update(ctx context.Context, cert *model.ConfigCertificate) (bool, error)

	// Delete 通过证书数据库 id 删除证书，并返回该证书删除时是否启用中。
	Delete(ctx context.Context, ids []int64) (bool, error)

	FindIDs(ctx context.Context, ids []int64) ([]*model.ConfigCertificate, error)

	Page(ctx context.Context, cond []gen.Condition, scope PageScope) (*Page[*model.ConfigCertificate], error)
}

func NewConfigCertificate(qry *query.Query) ConfigCertificate {
	return &configCertificateRepository{
		Repository: Base[*model.ConfigCertificate](qry),
		qry:        qry,
	}
}

type configCertificateRepository struct {
	Repository[*model.ConfigCertificate]
	qry   *query.Query
	mutex sync.Mutex
}

func (ccr *configCertificateRepository) Enables(ctx context.Context) ([]*model.ConfigCertificate, error) {
	tbl := ccr.qry.ConfigCertificate
	return tbl.WithContext(ctx).
		Where(tbl.Enabled.Is(true)).
		Find()
}

func (ccr *configCertificateRepository) Create(ctx context.Context, cert *model.ConfigCertificate, limit int64) (overflow, enabled bool, err error) {
	ccr.mutex.Lock()
	defer ccr.mutex.Unlock()

	err = ccr.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		dao := tbl.WithContext(ctx)
		if cnt, exx := dao.Count(); exx != nil {
			return exx
		} else if cnt >= limit {
			overflow = true
			return nil
		}

		return dao.Create(cert)
	})

	return
}

func (ccr *configCertificateRepository) Update(ctx context.Context, cert *model.ConfigCertificate) (bool, error) {
	ccr.mutex.Lock()
	defer ccr.mutex.Unlock()

	var enabled bool
	err := ccr.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		dao := tbl.WithContext(ctx)
		id := cert.ID
		dat, err := dao.Where(tbl.ID.Eq(id)).First()
		if err != nil {
			return err
		}
		enabled = dat.Enabled || cert.Enabled
		cert.CreatedAt = dat.CreatedAt

		return dao.Save(cert)
	})

	return enabled, err
}

func (ccr *configCertificateRepository) Delete(ctx context.Context, ids []int64) (bool, error) {
	ccr.mutex.Lock()
	defer ccr.mutex.Unlock()

	var enabled bool
	err := ccr.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		dao := tbl.WithContext(ctx)
		expr := tbl.ID.In(ids...)
		cnt, _ := dao.Where(expr, tbl.Enabled.Is(true)).Count()
		enabled = cnt > 0

		_, err := tbl.WithContext(ctx).Where(expr).Delete()

		return err
	})

	return enabled, err
}

func (ccr *configCertificateRepository) FindIDs(ctx context.Context, ids []int64) ([]*model.ConfigCertificate, error) {
	tbl := ccr.qry.ConfigCertificate
	return tbl.WithContext(ctx).Where(tbl.ID.In(ids...)).Find()
}

func (ccr *configCertificateRepository) Page(ctx context.Context, cond []gen.Condition, ps PageScope) (*Page[*model.ConfigCertificate], error) {
	tbl := ccr.qry.ConfigCertificate
	dao := tbl.WithContext(ctx).Where(cond...)
	count, err := dao.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return ccr.emptyRecords(ps), nil
	}

	dats, err := dao.Scopes(ps.Gen(count)).Find()
	if err != nil {
		return nil, err
	}

	return ccr.withRecords(ps, count, dats), nil
}
