package middle

import (
	"strings"

	"github.com/xgfone/ship/v5"
)

// NewAuth 需要凭借临时密钥才可以初始化。
func NewAuth(token string) ship.Middleware {
	m := stringToken(token)

	return m.handle
}

type stringToken string

func (st stringToken) handle(h ship.Handler) ship.Handler {
	return func(c *ship.Context) error {
		str := c.GetReqHeader(ship.HeaderAuthorization)
		if strings.EqualFold(str, string(st)) {
			return h(c)
		}

		return ship.ErrUnauthorized
	}
}
