package service

import (
	"log/slog"

	"github.com/dop251/goja"

	"github.com/xmx/aegis-server/jsenv/jsvm"
)

type Player interface {
	NewGoja(loads []jsvm.Loader) GojaPlayer
}

func NewPlayer(loads []jsvm.Loader, log *slog.Logger) Player {
	return &player{
		loads: loads,
		log:   log,
	}
}

type player struct {
	loads []jsvm.Loader
	log   *slog.Logger
}

func (pl *player) NewGoja(loads []jsvm.Loader) GojaPlayer {
	pugs := append(pl.loads, loads...)
	return NewGojaPlayer(pugs, pl.log)
}

type VM struct {
	vm *goja.Runtime
}

func NewVM(loads []jsvm.Loader) (*VM, error) {
	vm := jsvm.New()
	if err := jsvm.Register(vm, loads); err != nil {
		return nil, err
	}

	return &VM{vm: vm}, nil
}
