package gridfs

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

type FS interface {
	fs.FS
	OpenID(fileID int64) (File, error)
	Save(ctx context.Context, filename string, r io.Reader) (*model.GridFile, error)
}

type File interface {
	fs.File
	fs.FileInfo
	SHA1() string
	SHA256() string
}

func NewFS(qry *query.Query) FS {
	return &gridFS{
		qry:     qry,
		burst:   63 * 1024, // 63KiB
		timeout: time.Minute,
	}
}

type gridFS struct {
	qry     *query.Query
	burst   uint16
	timeout time.Duration
}

// Open 请确保
func (g *gridFS) Open(name string) (fs.File, error) {
	tbl := g.qry.GridFile
	ctx, cancel := g.getContext()
	defer cancel()

	file, err := tbl.WithContext(ctx).
		Where(tbl.Filename.Eq(name)).
		Order(tbl.CreatedAt.Desc()).
		First()
	if err != nil {
		return nil, err
	}

	f := g.newFile(file)

	return f, nil
}

func (g *gridFS) OpenID(fileID int64) (File, error) {
	tbl := g.qry.GridFile
	ctx, cancel := g.getContext()
	defer cancel()

	file, err := tbl.WithContext(ctx).
		Where(tbl.ID.Eq(fileID)).
		First()
	if err != nil {
		return nil, err
	}

	f := g.newFile(file)

	return f, nil
}

func (g *gridFS) Save(ctx context.Context, filename string, r io.Reader) (*model.GridFile, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	createdAt := time.Now()
	h1, h256 := sha1.New(), sha256.New()
	fr := io.TeeReader(r, io.MultiWriter(h1, h256))

	var sequence int64
	file := &model.GridFile{Filename: filename, Extension: ext, Burst: g.burst, CreatedAt: createdAt}
	err := g.qry.Transaction(func(tx *query.Query) error {
		ftbl, ctbl := tx.GridFile, tx.GridChunk
		fdao, cdao := ftbl.WithContext(ctx), ctbl.WithContext(ctx)

		// 先保存文件信息获取到数据库文件 ID。
		if err := fdao.Create(file); err != nil {
			return err
		}

		fileID := file.ID
		var length int64
		for {
			data := make([]byte, g.burst)
			n, err := io.ReadFull(fr, data)
			if n > 0 {
				length += int64(n)
				chunk := &model.GridChunk{FileID: fileID, Sequence: sequence, Data: data[:n]}
				sequence++
				if exx := cdao.Create(chunk); exx != nil {
					return exx
				}
			}

			if err != nil {
				// io.ReadFull returned error
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				return err
			}
		}

		// 更新文件信息
		h1sum, h256sum := h1.Sum(nil), h256.Sum(nil)
		file.SHA1 = hex.EncodeToString(h1sum)
		file.SHA256 = hex.EncodeToString(h256sum)
		file.Length = length
		file.UpdatedAt = time.Now()

		_, err := fdao.Updates(file)

		return err
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (g *gridFS) newFile(file *model.GridFile) File {
	return &gridFile{
		qry:     g.qry,
		file:    file,
		timeout: g.timeout,
	}
}

func (g *gridFS) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), g.timeout)
}
