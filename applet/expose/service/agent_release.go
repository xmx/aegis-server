package service

import (
	"context"
	"debug/buildinfo"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/applet/expose/request"
	"github.com/xmx/aegis-server/applet/expose/response"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (ar *AgentRelease) Parse(r io.ReaderAt) (*model.BuildInfo, *response.ExecutableMetadata, error) {
	bi, err := buildinfo.Read(r)
	if err != nil {
		return nil, nil, err
	}

	build := model.FormatBuildInfo(bi)
	info := new(response.ExecutableMetadata)
	for _, set := range bi.Settings {
		key, val := set.Key, set.Value
		switch key {
		case "GOOS":
			info.Goos = val
		case "GOARCH":
			info.Goarch = val
		case "vcs.time":
			if at, _ := time.ParseInLocation(time.RFC3339, val, time.UTC); !at.IsZero() {
				info.Version = ar.formatVersion(at)
			}
		}
	}
	if info.Version != "" {
		return build, info, nil
	}

	mv := bi.Main.Version
	after, _ := strings.CutPrefix(mv, "v0.0.0-")
	before, _, _ := strings.Cut(after, "-")
	at, _ := time.ParseInLocation("20060102150405", before, time.UTC)
	if !at.IsZero() {
		info.Version = ar.formatVersion(at)
	}

	return build, info, nil
}

func (ar *AgentRelease) Open(ctx context.Context, id bson.ObjectID) (*model.AgentRelease, *mongo.GridFSDownloadStream, error) {
	repo := ar.repo.AgentRelease()
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

func (ar *AgentRelease) formatVersion(t time.Time) string {
	return t.Format("06.1.2-150405")
}
