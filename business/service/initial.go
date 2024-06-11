package service

import (
	"context"
	"sync/atomic"

	"github.com/xmx/aegis-server/argument/bizdata"
)

type InitialService interface {
	Store(data *bizdata.InitialData) bool
	Wait(ctx context.Context) (*bizdata.InitialData, error)
}

func Initial() InitialService {
	wait := make(chan struct{})

	return &initialService{
		wait: wait,
	}
}

type initialService struct {
	wait chan struct{}
	done atomic.Bool
	data *bizdata.InitialData
}

func (svc *initialService) Store(data *bizdata.InitialData) bool {
	if svc.done.CompareAndSwap(false, true) {
		svc.data = data
		close(svc.wait)
		return true
	}
	return false
}

func (svc *initialService) Wait(ctx context.Context) (*bizdata.InitialData, error) {
	select {
	case <-svc.wait:
		return svc.data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
