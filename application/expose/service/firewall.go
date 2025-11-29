package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/application/expose/request"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewFirewall(repo repository.All, log *slog.Logger) *Firewall {
	return &Firewall{
		repo: repo,
		log:  log,
	}
}

type Firewall struct {
	repo repository.All
	log  *slog.Logger
}

func (fw *Firewall) Create(ctx context.Context, req *request.FirewallUpsert) error {
	now := time.Now()
	dat, err := req.Format()
	if err != nil {
		return err
	}

	mod := &model.Firewall{
		Name:         dat.Name,
		Enabled:      dat.Enabled,
		TrustHeaders: dat.TrustHeaders,
		TrustProxies: dat.TrustProxies,
		Blacklist:    dat.Blacklist,
		Inets:        dat.Inets,
		Countries:    dat.Inets,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	repo := fw.repo.Firewall()
	_, err = repo.InsertOne(ctx, mod)
	// TODO 清除缓存

	return err
}

func (fw *Firewall) Update(ctx context.Context, req *request.FirewallUpsert) error {
	now := time.Now()
	dat, err := req.Format()
	if err != nil {
		return err
	}

	mod := &model.Firewall{
		Name:         dat.Name,
		Enabled:      dat.Enabled,
		TrustHeaders: dat.TrustHeaders,
		TrustProxies: dat.TrustProxies,
		Blacklist:    dat.Blacklist,
		Inets:        dat.Inets,
		Countries:    dat.Countries,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	repo := fw.repo.Firewall()
	_, err = repo.InsertOne(ctx, mod)
	// TODO 清除缓存

	return err
}

func (fw *Firewall) Delete(ctx context.Context, name string) error {
	repo := fw.repo.Firewall()
	last, err := repo.FindOneAndDelete(ctx, bson.D{{"name", name}})
	if err != nil || last == nil {
		return err
	}

	// TODO 清除缓存

	return err
}
