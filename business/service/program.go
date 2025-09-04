package service

import (
	"context"
	"log/slog"
	"mime/multipart"

	"github.com/xmx/aegis-control/datalayer/repository"
)

func NewProgram(repo repository.All, log *slog.Logger) *Program {
	return &Program{repo: repo, log: log}
}

type Program struct {
	repo repository.All
	log  *slog.Logger
}

func (pgm *Program) Install(ctx context.Context, file *multipart.FileHeader) error {
	open, err := file.Open()
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer open.Close()

	return err
}
