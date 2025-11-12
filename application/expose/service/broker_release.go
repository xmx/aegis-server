package service

import (
	"context"
	"debug/buildinfo"
	"io"
	"log/slog"
	"time"

	"github.com/xmx/aegis-common/banner"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/application/errcode"
	"github.com/xmx/aegis-server/application/expose/request"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewBrokerRelease(repo repository.All, log *slog.Logger) *BrokerRelease {
	return &BrokerRelease{
		repo: repo,
		log:  log,
	}
}

type BrokerRelease struct {
	repo repository.All
	log  *slog.Logger
}

func (br *BrokerRelease) Upload(ctx context.Context, req *request.BrokerReleaseUpload) error {
	now := time.Now()
	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	bi, meta, err := br.Parse(file)
	if err != nil {
		return err
	}

	// 上传文件
	filename := req.File.Filename
	repo := br.repo.BrokerRelease()
	info, err := repo.SaveFile(ctx, file, filename)
	if err != nil {
		return err
	}

	build := model.FormatBuildInfo(bi)
	semver := model.ParseSemver(meta.Version)
	data := &model.BrokerRelease{
		FileID:    info.FileID,
		Filename:  filename,
		Goos:      meta.Goos,
		Goarch:    meta.Goarch,
		Length:    info.Length,
		Semver:    semver.Version,
		Version:   semver.Number,
		BuildInfo: build,
		Checksum:  info.Checksum,
		Changelog: req.Changelog,
		CreatedAt: now,
	}
	if _, err = repo.InsertOne(ctx, data); err == nil {
		return nil
	}
	_ = repo.DeleteFile(ctx, info.FileID)

	return err
}

func (br *BrokerRelease) Parse(r io.ReaderAt) (*buildinfo.BuildInfo, *banner.Info, error) {
	bi, err := buildinfo.Read(r)
	if err != nil {
		return nil, nil, err
	}

	if bi.Main.Path != "github.com/xmx/aegis-broker" {
		return nil, nil, errcode.ErrInvalidBinaryRelease
	}
	info := banner.ParseInfo(bi)

	return bi, info, nil
}

func (br *BrokerRelease) Open(ctx context.Context, id bson.ObjectID) (*model.AgentRelease, *mongo.GridFSDownloadStream, error) {
	repo := br.repo.AgentRelease()
	release, err := repo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	fileID := release.FileID
	stm, err := repo.OpenFile(ctx, fileID)
	if err != nil {
		return nil, nil, err
	}

	return release, stm, nil
}
