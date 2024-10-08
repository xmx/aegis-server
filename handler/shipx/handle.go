package shipx

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/library/validation"
	"gorm.io/gorm"
)

func NotFound(_ *ship.Context) error {
	return ship.ErrNotFound.Newf("资源不存在")
}

func HandleError(c *ship.Context, e error) {
	statusCode, title, detail := UnpackError(e)

	pd := &response.ProblemDetails{
		Title:    title,
		Status:   statusCode,
		Detail:   detail,
		Instance: c.Path(),
		Method:   c.Method(),
		Datetime: time.Now().UTC(),
	}
	_ = c.JSON(statusCode, pd)
}

func UnpackError(err error) (statusCode int, title string, detail string) {
	statusCode = http.StatusBadRequest
	title = "请求错误"
	detail = err.Error()

	switch ce := err.(type) {
	case ship.HTTPServerError:
		statusCode = ce.Code
	case *ship.HTTPServerError:
		statusCode = ce.Code
	case *validation.Error, *validation.NilError:
		title = "参数校验错误"
	case *time.ParseError:
		detail = "时间格式错误，正确格式：" + ce.Layout
	case *net.ParseError:
		detail = ce.Text + " 不是有效的 " + ce.Type
	case base64.CorruptInputError:
		detail = "Base64 编码错误：" + ce.Error()
	case *json.SyntaxError:
		detail = "错误的 JSON 格式"
	case *json.UnmarshalTypeError:
		detail = ce.Field + " 收到无效的数据类型"
	case *strconv.NumError:
		var msg string
		if sn := strings.SplitN(ce.Func, "Parse", 2); len(sn) == 2 {
			msg = ce.Num + " 不是 " + strings.ToLower(sn[1]) + " 类型"
		} else {
			msg = "类型错误：" + ce.Num
		}
		detail = msg
	case *http.MaxBytesError:
		limit := strconv.FormatInt(ce.Limit, 10)
		detail = "请求报文超过 " + limit + " 个字节限制"
		statusCode = http.StatusRequestEntityTooLarge
	default:
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			detail = "数据不存在"
		case errors.Is(err, ship.ErrSessionNotExist), errors.Is(err, ship.ErrInvalidSession):
			statusCode = http.StatusUnauthorized
			detail = "认证无效"
		}
	}

	return
}
