package middle

import (
	"strings"

	"github.com/labstack/echo/v5"
)

// NewAuth 需要凭借临时密钥才可以初始化。
func NewAuth(token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			str := c.Request().Header.Get(echo.HeaderAuthorization)
			if strings.EqualFold(str, token) {
				return next(c)
			}

			return echo.ErrUnauthorized
		}
	}
}
