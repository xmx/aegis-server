package model

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type CertificateConfigurer interface {
	Create(ctx context.Context, cert *model.Certificate) error
	Update(ctx context.Context, cert *model.Certificate) error
	Delete(ctx context.Context, id int64) error
	Effect(ctx context.Context) (*model.Certificate, error)
}

func Certificate(qry *query.Query) CertificateConfigurer {
	return &certificateConfig{
		qry: qry,
	}
}

type certificateConfig struct {
	qry *query.Query
}

func (c *certificateConfig) Create(ctx context.Context, cert *model.Certificate) error {
	enabled := cert.Enabled
	if !enabled {
		tbl := c.qry.Certificate
		return tbl.WithContext(ctx).Create(cert)
	}

	return c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.Certificate
		if _, err := tbl.WithContext(ctx).
			Where(tbl.Enabled.Is(true)).
			Update(tbl.Enabled, false); err != nil {
			return err
		}
		return tbl.WithContext(ctx).Create(cert)
	})
}

func (c *certificateConfig) Update(ctx context.Context, cert *model.Certificate) error {
	id, enabled := cert.ID, cert.Enabled
	if !enabled {
		tbl := c.qry.Certificate
		_, err := tbl.WithContext(ctx).
			Where(tbl.ID.Eq(id)).
			UpdateColumns(cert)
		return err
	}

	return c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.Certificate
		if _, err := tbl.WithContext(ctx).
			Where(tbl.Enabled.Is(true)).
			Update(tbl.Enabled, false); err != nil {
			return err
		}
		_, err := tbl.WithContext(ctx).
			Where(tbl.ID.Eq(id)).
			UpdateColumns(cert)
		return err
	})
}

func (c *certificateConfig) Delete(ctx context.Context, id int64) error {
	tbl := c.qry.Certificate
	_, err := tbl.WithContext(ctx).
		Where(tbl.ID.Value(id)).
		Delete()
	return err
}

func (c *certificateConfig) Effect(ctx context.Context) (*model.Certificate, error) {
	tbl := c.qry.Certificate
	return tbl.WithContext(ctx).
		Where(tbl.Enabled.Is(true)).
		First()
}
