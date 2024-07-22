package errcode

import "github.com/xgfone/ship/v5"

var ErrUnsupportedHijack = ship.ErrBadRequest.Newf("不支持连接升级")
