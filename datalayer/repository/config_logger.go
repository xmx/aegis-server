package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type ConfigLoggerRepository interface {
	Enabled(ctx context.Context) (*model.ConfigLogger, error)
	Create(ctx context.Context, cert *model.ConfigLogger) (enabled bool, err error)
	Delete(ctx context.Context, id int64) (enabled bool, err error)
}

func ConfigLogger(qry *query.Query) ConfigLoggerRepository {
	return &configLoggerRepository{qry: qry}
}

type configLoggerRepository struct {
	qry *query.Query
}

func (c *configLoggerRepository) Enabled(ctx context.Context) (*model.ConfigLogger, error) {
	tbl := c.qry.ConfigLogger
	return tbl.WithContext(ctx).
		Where(tbl.Enabled.Is(true)).
		First()
}

func (c *configLoggerRepository) Create(ctx context.Context, cert *model.ConfigLogger) (bool, error) {
	enabled := cert.Enabled
	tbl := c.qry.ConfigLogger
	err := tbl.WithContext(ctx).Create(cert)

	return enabled, err
}

func (c *configLoggerRepository) Delete(ctx context.Context, id int64) (bool, error) {
	var enabled bool
	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigLogger
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
