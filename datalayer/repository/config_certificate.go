package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type ConfigCertificateRepository interface {
	Enabled(ctx context.Context) (*model.ConfigCertificate, error)
	Create(ctx context.Context, cert *model.ConfigCertificate) (enabled bool, err error)
	Update(ctx context.Context, cert *model.ConfigCertificate) (enabled bool, err error)
	Delete(ctx context.Context, id int64) (enabled bool, err error)
}

func ConfigCertificate(qry *query.Query) ConfigCertificateRepository {
	return &configCertificateRepository{qry: qry}
}

type configCertificateRepository struct {
	qry *query.Query
}

func (c *configCertificateRepository) Enabled(ctx context.Context) (*model.ConfigCertificate, error) {
	tbl := c.qry.ConfigCertificate
	return tbl.WithContext(ctx).
		Where(tbl.Enabled.Is(true)).
		First()
}

func (c *configCertificateRepository) Create(ctx context.Context, cert *model.ConfigCertificate) (bool, error) {
	enabled := cert.Enabled
	if !enabled {
		tbl := c.qry.ConfigCertificate
		err := tbl.WithContext(ctx).Create(cert)
		return false, err
	}

	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		if _, err := tbl.WithContext(ctx).
			Where(tbl.Enabled.Is(true)).
			UpdateSimple(tbl.Enabled.Value(false)); err != nil {
			return err
		}
		return tbl.WithContext(ctx).Create(cert)
	})

	return enabled, err
}

func (c *configCertificateRepository) Update(ctx context.Context, cert *model.ConfigCertificate) (bool, error) {
	enabled := cert.Enabled
	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		dao := tbl.WithContext(ctx)
		id := cert.ID
		dat, err := dao.Where(tbl.ID.Eq(id)).First()
		if err != nil {
			return err
		}

		if enabled = enabled || enabled != dat.Enabled; enabled {
			if _, err = dao.Where(tbl.Enabled.Is(true)).
				UpdateSimple(tbl.Enabled.Value(false)); err != nil {
				return err
			}
		}
		cert.UpdatedAt = dat.UpdatedAt
		_, err = dao.Updates(cert)

		return err
	})

	return enabled, err
}

func (c *configCertificateRepository) Delete(ctx context.Context, id int64) (bool, error) {
	var enabled bool
	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigCertificate
		dat, err := tbl.WithContext(ctx).Where(tbl.ID.Eq(id)).First()
		if err != nil {
			return err
		}

		enabled = dat.Enabled
		_, err = tbl.WithContext(ctx).Where(tbl.ID.Eq(id)).Delete()

		return err
	})

	return enabled, err
}
