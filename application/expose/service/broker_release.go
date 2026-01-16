package service

import (
	"cmp"
	"context"
	"debug/buildinfo"
	"io"
	"log/slog"
	"slices"
	"time"

	"github.com/xmx/aegis-common/banner"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/application/errcode"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/response"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (br *BrokerRelease) List(ctx context.Context) ([]*model.BrokerRelease, error) {
	order := bson.D{{"version", -1}, {"goos", 1}, {"goarch", 1}}
	opt := options.Find().SetSort(order)
	repo := br.repo.BrokerRelease()

	return repo.Find(ctx, bson.D{}, opt)
}

func (br *BrokerRelease) Delete(ctx context.Context, id bson.ObjectID) error {
	repo := br.repo.BrokerRelease()
	dat, err := repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err = repo.DeleteFile(ctx, dat.FileID); err != nil {
		return err
	}
	_, err = repo.DeleteByID(ctx, id)

	return err
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
	semver := model.ParseSemver(meta.Semver)
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

func (br *BrokerRelease) Open(ctx context.Context, fileID bson.ObjectID) (*mongo.GridFSDownloadStream, error) {
	repo := br.repo.BrokerRelease()
	stm, err := repo.OpenFile(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return stm, nil
}

func (br *BrokerRelease) Exposes(ctx context.Context) (model.ExposeAddresses, error) {
	repo := br.repo.Setting()
	dat, err := repo.Get(ctx)
	if err != nil {
		return nil, err
	} else if len(dat.Exposes) == 0 {
		return nil, errcode.ErrNilDocument
	}

	return dat.Exposes, nil
}

func (br *BrokerRelease) Latest(ctx context.Context, goos, goarch string) (*model.BrokerRelease, error) {
	filter := bson.D{{"goos", goos}, {"goarch", goarch}}
	order := bson.D{{"version", -1}, {"_id", -1}}
	opt := options.FindOne().SetSort(order)
	repo := br.repo.BrokerRelease()

	return repo.FindOne(ctx, filter, opt)
}

func (br *BrokerRelease) Platforms(ctx context.Context) ([]*response.FieldValues[string, string], error) {
	pipe := mongo.Pipeline{
		{{"$group", bson.D{{"_id", "$goos"}, {"values", bson.D{{"$addToSet", "$goarch"}}}}}},
		{{"$project", bson.D{{"field", "$_id"}, {"values", 1}, {"_id", 0}}}},
	}

	var ret []*response.FieldValues[string, string]
	repo := br.repo.BrokerRelease()
	if err := repo.AggregateTo(ctx, pipe, &ret); err != nil {
		return nil, err
	}
	slices.SortFunc(ret, func(a, b *response.FieldValues[string, string]) int {
		return cmp.Compare(a.Field, b.Field)
	})
	for _, r := range ret {
		slices.SortFunc(r.Values, func(a, b string) int {
			return cmp.Compare(a, b)
		})
	}

	return ret, nil
}
