package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/library/memoize"
	"github.com/xmx/aegis-server/application/expose/firewalld"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/library/iplist"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewFirewall(repo repository.All, mmr firewalld.MaxmindReader, log *slog.Logger) *Firewall {
	fw := &Firewall{
		repo: repo,
		mmr:  mmr,
		log:  log,
	}
	fw.cfg = memoize.NewCache2(fw.slowLoad)

	return fw
}

type Firewall struct {
	repo repository.All
	mmr  firewalld.MaxmindReader
	log  *slog.Logger
	cfg  memoize.Cache2[*firewalld.Config, error]
}

func (fw *Firewall) List(ctx context.Context) ([]*model.Firewall, error) {
	return fw.repo.Firewall().Find(ctx, bson.D{})
}

func (fw *Firewall) Create(ctx context.Context, req *request.FirewallUpsert) error {
	now := time.Now()
	dat, err := req.Format()
	if err != nil {
		return err
	}

	enabled := dat.Enabled
	mod := &model.Firewall{
		Name:         dat.Name,
		Enabled:      enabled,
		TrustProxies: dat.TrustProxies,
		TrustHeaders: dat.TrustHeaders,
		Blacklist:    dat.Blacklist,
		CountryMode:  dat.CountryMode,
		IPNets:       dat.IPNets,
		Countries:    dat.Countries,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	repo := fw.repo.Firewall()
	if _, err = repo.InsertOne(ctx, mod); err != nil || !enabled {
		return err
	}
	fw.Reset()

	return nil
}

func (fw *Firewall) Update(ctx context.Context, req *request.FirewallUpsert) error {
	now := time.Now()
	dat, err := req.Format()
	if err != nil {
		return err
	}

	name, enabled := req.Name, req.Enabled
	mod := &model.Firewall{
		Enabled:      dat.Enabled,
		TrustHeaders: dat.TrustHeaders,
		TrustProxies: dat.TrustProxies,
		Blacklist:    dat.Blacklist,
		CountryMode:  dat.CountryMode,
		IPNets:       dat.IPNets,
		Countries:    dat.Countries,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	repo := fw.repo.Firewall()
	filter := bson.D{{"name", name}}
	update := bson.M{"$set": mod}
	last, err := repo.FindOneAndUpdate(ctx, filter, update)
	if err != nil {
		return err
	}
	if enabled || enabled != last.Enabled {
		fw.Reset()
	}

	return err
}

func (fw *Firewall) Delete(ctx context.Context, name string) error {
	repo := fw.repo.Firewall()
	last, err := repo.FindOneAndDelete(ctx, bson.D{{"name", name}})
	if err != nil || last == nil || !last.Enabled {
		return err
	}
	fw.Reset()

	return err
}

func (fw *Firewall) Configure(ctx context.Context) (*firewalld.Config, error) {
	return fw.cfg.Load(ctx)
}

func (fw *Firewall) Reset() {
	_, _ = fw.cfg.Forget()
}

func (fw *Firewall) Sandbox(req *request.FirewallUpsert) (*firewalld.Firewalld, error) {
	dat, err := req.Format()
	if err != nil {
		return nil, err
	}
	cf := func(context.Context) (*firewalld.Config, error) {
		return fw.parseConfig(dat)
	}
	cfg := firewalld.ConfigureFunc(cf)
	wall := firewalld.New(cfg, fw.log)

	return wall, nil
}

func (fw *Firewall) slowLoad(ctx context.Context) (*firewalld.Config, error) {
	mod, err := fw.repo.Firewall().Enabled(ctx)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	dat := request.FirewallUpsert{
		Name:         mod.Name,
		Enabled:      mod.Enabled,
		Blacklist:    mod.Blacklist,
		TrustHeaders: mod.TrustHeaders,
		TrustProxies: mod.TrustProxies,
		CountryMode:  mod.CountryMode,
		IPNets:       mod.IPNets,
		Countries:    mod.Countries,
	}

	return fw.parseConfig(dat)
}

func (fw *Firewall) parseConfig(dat request.FirewallUpsert) (*firewalld.Config, error) {
	cfg := &firewalld.Config{
		TrustHeaders: dat.TrustHeaders,
		Blacklist:    dat.Blacklist,
		CountryMode:  dat.CountryMode,
		MaxmindDB:    fw.mmr,
	}
	if len(dat.TrustProxies) != 0 {
		proxies, err := iplist.Parse(dat.TrustProxies)
		if err != nil {
			return nil, err
		}
		cfg.TrustProxies = proxies
	}
	if !dat.CountryMode && len(dat.IPNets) != 0 {
		inets, err := iplist.Parse(dat.IPNets)
		if err != nil {
			return nil, err
		}
		cfg.IPNets = inets
	}
	if dat.CountryMode && len(dat.Countries) != 0 {
		countries := make(map[string]struct{}, len(dat.Countries))
		for _, country := range dat.Countries {
			countries[country] = struct{}{}
		}
		cfg.Countries = countries
	}

	return cfg, nil
}
