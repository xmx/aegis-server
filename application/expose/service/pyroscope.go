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

func NewPyroscope(repo repository.All, log *slog.Logger) *Pyroscope {
	return &Pyroscope{
		repo: repo,
		log:  log,
	}
}

type Pyroscope struct {
	repo repository.All
	log  *slog.Logger
}

func (prs *Pyroscope) List(ctx context.Context) ([]*model.Pyroscope, error) {
	repo := prs.repo.Pyroscope()
	return repo.Find(ctx, bson.D{})
}

func (prs *Pyroscope) Create(ctx context.Context, req *request.PyroscopeUpsert) error {
	now := time.Now()
	name, enabled := req.Name, req.Enabled
	repo := prs.repo.Pyroscope()
	data := &model.Pyroscope{
		Name:      name,
		Address:   req.Address,
		Username:  req.Username,
		Password:  req.Password,
		Enabled:   enabled,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := repo.InsertOne(ctx, data); err != nil || !enabled {
		return err
	}

	// 自动关闭其他开启的配置
	filter := bson.D{{"enabled", true}, {"name", bson.D{{"$ne", name}}}}
	update := bson.M{"$set": bson.M{"enabled": false}}
	if _, err := repo.UpdateMany(ctx, filter, update); err != nil {
		prs.log.ErrorContext(ctx, "自动关闭其他启用的配置发生错误", "error", err)
		return err
	}

	//prs.reset()

	return nil
}

func (prs *Pyroscope) Update(ctx context.Context, req *request.PyroscopeUpsert) error {
	now := time.Now()
	name, enabled := req.Name, req.Enabled
	repo := prs.repo.Pyroscope()
	filter := bson.D{{"name", name}}
	data := &model.Pyroscope{
		Address:   req.Address,
		Enabled:   enabled,
		UpdatedAt: now,
	}
	update := bson.M{"$set": data}

	last, err := repo.FindOneAndUpdate(ctx, filter, update)
	if err != nil {
		return err
	}
	if enabled {
		filter = bson.D{{"enabled", true}, {"name", bson.D{{"$ne", name}}}}
		update = bson.M{"$set": bson.M{"enabled": false}}
		if _, err = repo.UpdateMany(ctx, filter, update); err != nil {
			prs.log.ErrorContext(ctx, "自动关闭其他启用的配置发生错误", "error", err)
			return err
		}
	}
	if last == nil || enabled || enabled != last.Enabled {
		//prs.reset()
	}

	return nil
}

func (prs *Pyroscope) Delete(ctx context.Context, name string) error {
	repo := prs.repo.Pyroscope()
	filter := bson.D{{"name", name}}

	last, err := repo.FindOneAndDelete(ctx, filter)
	if err != nil || last == nil || !last.Enabled {
		return err
	}
	//prs.reset()

	return nil
}
