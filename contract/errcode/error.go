package errcode

import "github.com/xmx/ship"

var (
	ErrDataNotExists          = ship.ErrNotFound.Newf("数据不存在")
	ErrCertificateUnavailable = ship.ErrBadRequest.Newf("未配置有效的证书")
	ErrCertificateInvalid     = ship.ErrBadRequest.Newf("无效证书")
)
