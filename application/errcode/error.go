package errcode

import (
	"github.com/xgfone/ship/v5"
)

var (
	ErrForbiddenOAuthProvider   = ship.ErrForbidden.Newf("当前认证方式已暂停使用")
	ErrUnsupportedOAuthProvider = ship.ErrForbidden.Newf("不支持该认证方式")
	ErrOOBEFinished             = ship.ErrForbidden.Newf("项目已经初始化完毕")
)
