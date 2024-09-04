package gridfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"time"

	"gorm.io/gorm"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

func newFile(qry *query.Query, file *model.GridFile) *gridFile {
	return &gridFile{
		qry:  qry,
		file: file,
	}
}

type gridFile struct {
	qry      *query.Query
	file     *model.GridFile
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
	for n < size {
		num := copy(b[n:], g.data[g.cursor:])
		if num == 0 {
			g.err = g.readNext()
		}
		if err := g.err; err != nil {
			if n > 0 {
				return n, nil
			} else {
				return 0, err
			}
		}
		n += num
		g.cursor += num
	}

	return n, nil
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

func (g *gridFile) SHA1() string {
	return g.file.SHA1
}

func (g *gridFile) SHA256() string {
	return g.file.SHA256
}

func (g *gridFile) readNext() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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
