package service

import (
	"context"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

func NewConfigServer(qry *query.Query) *ConfigServer {
	return &ConfigServer{qry: qry}
}

type ConfigServer struct {
	qry *query.Query
}

func (c *ConfigServer) Enabled(ctx context.Context) (*model.ConfigServer, error) {
	tbl := c.qry.ConfigServer
	return tbl.WithContext(ctx).
		Where(tbl.Enabled.Is(true)).
		First()
}

func (c *ConfigServer) Create(ctx context.Context, cert *model.ConfigServer) (bool, error) {
	enabled := cert.Enabled
	tbl := c.qry.ConfigServer
	err := tbl.WithContext(ctx).Create(cert)

	return enabled, err
}

func (c *ConfigServer) Delete(ctx context.Context, id int64) (bool, error) {
	var enabled bool
	err := c.qry.Transaction(func(tx *query.Query) error {
		tbl := tx.ConfigServer
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
