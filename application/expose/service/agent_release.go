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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewAgentRelease(repo repository.All, log *slog.Logger) *AgentRelease {
	return &AgentRelease{
		repo: repo,
		log:  log,
	}
}

type AgentRelease struct {
	repo repository.All
	log  *slog.Logger
}

func (ar *AgentRelease) Upload(ctx context.Context, req *request.AgentReleaseUpload) error {
	now := time.Now()
	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	build, meta, err := ar.Parse(file)
	if err != nil {
		return err
	}

	// 上传文件
	filename := req.File.Filename
	repo := ar.repo.AgentRelease()
	info, err := repo.SaveFile(ctx, file, filename)
	if err != nil {
		return err
	}

	semver := model.ParseSemver(meta.Version)
	data := &model.AgentRelease{
		FileID:    info.FileID,
		Filename:  filename,
		Goos:      meta.Goos,
		Goarch:    meta.Goarch,
		Length:    info.Length,
		Semver:    semver.Version,
		Version:   semver.Number,
		BuildInfo: model.FormatBuildInfo(build),
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

func (ar *AgentRelease) Parse(r io.ReaderAt) (*buildinfo.BuildInfo, *banner.Info, error) {
	bi, err := buildinfo.Read(r)
	if err != nil {
		return nil, nil, err
	}
	if bi.Main.Path != "github.com/xmx/aegis-agent" {
		return nil, nil, errcode.ErrInvalidBinaryRelease
	}
	info := banner.ParseInfo(bi)

	return bi, info, nil
}

func (ar *AgentRelease) Open(ctx context.Context, fileID bson.ObjectID) (*mongo.GridFSDownloadStream, error) {
	repo := ar.repo.AgentRelease()
	stm, err := repo.OpenFile(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return stm, nil
}

func (ar *AgentRelease) Latest(ctx context.Context, goos, goarch string) (*model.AgentRelease, error) {
	filter := bson.D{{"goos", goos}, {"goarch", goarch}}
	order := bson.D{{"version", -1}, {"_id", -1}}
	opt := options.FindOne().SetSort(order)
	repo := ar.repo.AgentRelease()

	return repo.FindOne(ctx, filter, opt)
}

func (ar *AgentRelease) Exposes(ctx context.Context) (model.ExposeAddresses, error) {
	opt := options.Find().SetProjection(bson.M{"exposes": 1})
	repo := ar.repo.Broker()
	broks, err := repo.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	} else if len(broks) == 0 {
		return nil, errcode.ErrNilDocument
	}

	var exposes model.ExposeAddresses
	for _, brok := range broks {
		exposes = append(exposes, brok.Exposes...)
	}

	return exposes, nil
}
