package gridfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"time"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"gorm.io/gorm"
)

type gridFile struct {
	qry      *query.Query
	file     *model.GridFile
	timeout  time.Duration
	sequence int64
	data     []byte
	cursor   int
	err      error
}

func (g *gridFile) Stat() (fs.FileInfo, error) {
	return g, nil
}

func (g *gridFile) Read(b []byte) (int, error) {
	if err := g.err; err != nil {
		return 0, err
	}

	size := len(b)
	var n int
	for n < size && g.err == nil {
		num := copy(b[n:], g.data[g.cursor:])
		if num > 0 {
			n += num
			g.cursor += num
		} else {
			g.err = g.readNext()
		}
	}
	if n > 0 {
		return n, nil
	}

	return 0, g.err
}

func (g *gridFile) Close() error {
	return nil
}

func (g *gridFile) Name() string {
	return g.file.Filename
}

func (g *gridFile) Size() int64 {
	return g.file.Length
}

func (g *gridFile) Mode() fs.FileMode {
	return 0o644
}

func (g *gridFile) ModTime() time.Time {
	return g.file.UpdatedAt
}

func (g *gridFile) IsDir() bool {
	return false
}

// Sys underlying data source (can return nil)
func (g *gridFile) Sys() any {
	return nil
}

func (g *gridFile) MD5() string {
	return g.file.MD5
}

func (g *gridFile) SHA1() string {
	return g.file.SHA1
}

func (g *gridFile) SHA256() string {
	return g.file.SHA256
}

func (g *gridFile) readNext() error {
	ctx, cancel := g.getContext()
	defer cancel()

	tbl := g.qry.GridChunk
	chunk, err := tbl.WithContext(ctx).
		Where(tbl.FileID.Eq(g.file.ID), tbl.Sequence.Eq(g.sequence)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return io.EOF
		}
		return err
	}

	g.sequence++
	g.data = chunk.Data
	g.cursor = 0

	return nil
}

func (g *gridFile) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), g.timeout)
}
