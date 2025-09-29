package service

import (
	"context"
	"io/fs"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
)

type FS struct {
	repo repository.All
	log  *slog.Logger
}

func NewFS(repo repository.All, log *slog.Logger) *FS {
	return &FS{repo: repo, log: log}
}

func (fs *FS) Upload(ctx context.Context, dir string, mh *multipart.FileHeader) (*model.FS, error) {
	file, err := mh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	filename := mh.Filename
	fp := path.Join(dir, filename)

	repo := fs.repo.FS()

	return repo.Create(ctx, fp, file)
}

func (fs *FS) Mkdir(ctx context.Context, dir string) error {
	repo := fs.repo.FS()
	return repo.Mkdir(ctx, dir)
}

func (fs *FS) Entries(ctx context.Context, dir string) (model.FSs, error) {
	repo := fs.repo.FS()
	return repo.Entries(ctx, dir)
}

func (fs *FS) Open(ctx context.Context, dir string) (fs.File, error) {
	repo := fs.repo.FS()
	return repo.OpenContext(ctx, dir)
}

func (fs *FS) Remove(ctx context.Context, dir string) error {
	repo := fs.repo.FS()
	return repo.Remove(ctx, dir)
}

func (fs *FS) Handler() http.Handler {
	repo := fs.repo.FS()
	hfs := http.FS(repo)

	return http.FileServer(hfs)
}
