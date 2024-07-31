package errcode

import "github.com/xgfone/ship/v5"

var ErrConnectionHijack = ship.ErrBadRequest.Newf("不支持连接升级")

const (
	FmtTooManyCertificate = formatError("证书超过 %d 限制")
)

type formatError string

func (f formatError) Fmt(v ...any) error {
	return ship.ErrBadRequest.Newf(string(f), v...)
}
