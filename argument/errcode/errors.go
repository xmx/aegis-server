package errcode

import "github.com/xgfone/ship/v5"

var ErrConnectionHijack = ship.ErrBadRequest.Newf("不支持连接升级")
