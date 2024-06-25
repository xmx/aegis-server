package repository

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type CertificateRepository interface {
	Enabled(ctx context.Context) (*model.Certificate, error)
	Create(ctx context.Context, cert *model.Certificate) (enabled bool, err error)
	Delete(ctx context.Context, id int64) (enabled bool, err error)
}

func Certificate(qry *query.Query) CertificateRepository {
	return &certificateRepository{qry: qry}
}

type certificateRepository struct {
	qry *query.Query
}

func (c *certificateRepository) Enabled(ctx context.Context) (*model.Certificate, error) {
	tbl := c.qry.Certificate
	return tbl.WithContext(ctx).
		Where(tbl.Enabled.Is(true)).
		First()
}

func (c *certificateRepository) Create(ctx context.Context, cert *model.Certificate) (bool, error) {
	enabled := cert.Enabled
	tbl := c.qry.Certificate
	err := tbl.WithContext(ctx).Create(cert)

	return enabled, err
}

func (c *certificateRepository) Delete(ctx context.Context, id int64) (bool, error) {
	var enabled bool
	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.Certificate
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
