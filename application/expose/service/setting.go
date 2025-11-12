package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
)

func NewSetting(repo repository.All, log *slog.Logger) *Setting {
	return &Setting{repo: repo, log: log}
}

type Setting struct {
	repo repository.All
	log  *slog.Logger
}

func (st *Setting) Get(ctx context.Context) (*model.Setting, error) {
	repo := st.repo.Setting()
	return repo.Get(ctx)
}
