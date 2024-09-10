package errcode

import "github.com/xgfone/ship/v5"

var (
	ErrCertificateExisted = ship.ErrBadRequest.Newf("该证书已存在")
	ErrCertificateInvalid = ship.ErrBadRequest.Newf("无效的证书")
	ErrServerSentEvents   = ship.ErrBadRequest.Newf("请使用 Server-Sent Events 请求")
	ErrRequiredColumn     = ship.ErrBadRequest.Newf("字段名必须填写")

	ErrTooManyRequests = ship.ErrTooManyRequests.Newf("请求过多请稍后重试")
)

const (
	FmtTooManyCertificate = formatError("证书超过 %d 限制")
)

type formatError string

func (f formatError) Fmt(v ...any) error {
	return ship.ErrBadRequest.Newf(string(f), v...)
}
