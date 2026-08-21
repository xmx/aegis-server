package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
)

type OAuth struct {
	db  *repository.BaseDB
	log *slog.Logger
}

func NewOAuth(db *repository.BaseDB, log *slog.Logger) *OAuth {
	return &OAuth{
		db:  db,
		log: log,
	}
}

func (o *OAuth) Provider(ctx context.Context, provider string) (*model.OAuthClient, error) {
	coll := o.db.OAuthClient()
	return coll.FindByProvider(ctx, provider)
}
