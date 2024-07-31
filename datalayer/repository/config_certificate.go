package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type ConfigCertificateRepository interface {
	Repository[*model.ConfigCertificate]

	// Enables 查询已经启用的证书。
	Enables(ctx context.Context) ([]*model.ConfigCertificate, error)

	// Update 更新证书内容，用 ID 作为搜索条件。
	// 返回修改之前或修改之后是否启用。
	Update(ctx context.Context, cert *model.ConfigCertificate) (bool, error)

	// Delete 通过证书数据库 id 删除证书，并返回该证书删除时是否启用中。
	Delete(ctx context.Context, id int64) (bool, error)
}

func ConfigCertificate(qry *query.Query) ConfigCertificateRepository {
	return &configCertificateRepository{
		Repository: Base[*model.ConfigCertificate](qry),
		qry:        qry,
	}
}

type configCertificateRepository struct {
	Repository[*model.ConfigCertificate]
	qry *query.Query
}

func (c *configCertificateRepository) Enables(ctx context.Context) ([]*model.ConfigCertificate, error) {
	tbl := c.qry.ConfigCertificate
	return tbl.WithContext(ctx).
		Where(tbl.Enabled.Is(true)).
		Find()
}

func (c *configCertificateRepository) Create(ctx context.Context, cert *model.ConfigCertificate) (bool, error) {
	enabled := cert.Enabled
	tbl := c.qry.ConfigCertificate
	err := tbl.WithContext(ctx).Create(cert)

	return enabled, err
}

func (c *configCertificateRepository) Update(ctx context.Context, cert *model.ConfigCertificate) (bool, error) {
	var enabled bool
	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		dao := tbl.WithContext(ctx)
		id := cert.ID
		dat, err := dao.Where(tbl.ID.Eq(id)).First()
		if err != nil {
			return err
		}
		enabled = dat.Enabled || cert.Enabled
		cert.UpdatedAt = dat.UpdatedAt

		return dao.Save(cert)
	})

	return enabled, err
}

func (c *configCertificateRepository) Delete(ctx context.Context, id int64) (bool, error) {
	var enabled bool
	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		expr := tbl.ID.Eq(id)
		dat, err := tbl.WithContext(ctx).Where(expr).First()
		if err != nil {
			return err
		}

		enabled = dat.Enabled
		_, err = tbl.WithContext(ctx).Where(expr).Delete()

		return err
	})

	return enabled, err
}
