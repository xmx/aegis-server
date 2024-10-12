package service

import (
	"log/slog"

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

func NewVM() {
}
