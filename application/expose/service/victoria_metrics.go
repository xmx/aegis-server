package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/application/expose/request"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewVictoriaMetrics(repo repository.All, log *slog.Logger) *VictoriaMetrics {
	return &VictoriaMetrics{
		repo: repo,
		log:  log,
	}
}

type VictoriaMetrics struct {
	repo repository.All
	log  *slog.Logger
	mtx  sync.Mutex
}

func (vm *VictoriaMetrics) List(ctx context.Context) ([]*model.VictoriaMetrics, error) {
	repo := vm.repo.VictoriaMetrics()
	return repo.Find(ctx, bson.D{})
}

func (vm *VictoriaMetrics) Create(ctx context.Context, req *request.VictoriaMetricsUpsert) error {
	now := time.Now()
	name, enabled := req.Name, req.Enabled
	repo := vm.repo.VictoriaMetrics()
	data := &model.VictoriaMetrics{
		Name:      name,
		Method:    req.Method,
		Address:   req.Address,
		Header:    req.Header,
		Enabled:   enabled,
		UpdatedAt: now,
		CreatedAt: now,
	}

	vm.mtx.Lock()
	defer vm.mtx.Unlock()

	if _, err := repo.InsertOne(ctx, data); err != nil || !enabled {
		return err
	}

	// 自动关闭其他开启的配置
	filter := bson.D{{"enabled", true}, {"name", bson.D{{"$ne", name}}}}
	update := bson.M{"$set": bson.M{"enabled": false}}
	if _, err := repo.UpdateMany(ctx, filter, update); err != nil {
		vm.log.ErrorContext(ctx, "自动关闭其他启用的配置发生错误", "error", err)
		return err
	}

	vm.reset()

	return nil
}

func (vm *VictoriaMetrics) Update(ctx context.Context, req *request.VictoriaMetricsUpsert) error {
	now := time.Now()
	name, enabled := req.Name, req.Enabled
	repo := vm.repo.VictoriaMetrics()
	filter := bson.D{{"name", name}}
	data := &model.VictoriaMetrics{
		Method:    req.Method,
		Address:   req.Address,
		Header:    req.Header,
		Enabled:   enabled,
		UpdatedAt: now,
	}
	update := bson.M{"$set": data}

	vm.mtx.Lock()
	defer vm.mtx.Unlock()

	last, err := repo.FindOneAndUpdate(ctx, filter, update)
	if err != nil {
		return err
	}
	if enabled {
		filter = bson.D{{"enabled", true}, {"name", bson.D{{"$ne", name}}}}
		update = bson.M{"$set": bson.M{"enabled": false}}
		if _, err = repo.UpdateMany(ctx, filter, update); err != nil {
			vm.log.ErrorContext(ctx, "自动关闭其他启用的配置发生错误", "error", err)
			return err
		}
	}
	if last == nil || enabled || enabled != last.Enabled {
		vm.reset()
	}

	return nil
}

func (vm *VictoriaMetrics) Delete(ctx context.Context, name string) error {
	repo := vm.repo.VictoriaMetrics()
	filter := bson.D{{"name", name}}

	vm.mtx.Lock()
	defer vm.mtx.Unlock()
	last, err := repo.FindOneAndDelete(ctx, filter)
	if err != nil || last == nil || !last.Enabled {
		return err
	}
	vm.reset()

	return nil
}

func (vm *VictoriaMetrics) Reset() {
	vm.mtx.Lock()
	defer vm.mtx.Unlock()
	vm.reset()
}

func (vm *VictoriaMetrics) reset() {

}
