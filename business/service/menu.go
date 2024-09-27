package service

import (
	"context"
	"log/slog"

	"github.com/xgfone/ship/v5"
)

type Menu interface {
	Migrate(ctx context.Context, routes []ship.Route) error
}

type menuService struct {
	log *slog.Logger
}

func (svc menuService) Migrate(ctx context.Context, routes []ship.Route) error {
	for _, route := range routes {
		data := route.Data
		if data == nil {
			continue
		}
	}

	return nil
}
