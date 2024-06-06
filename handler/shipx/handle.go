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
)

func NotFound(_ *ship.Context) error {
	return ship.ErrNotFound
}

func HandleError(c *ship.Context, e error) {
	pd := &response.ProblemDetails{
		Type:     c.Host(),
		Title:    "请求错误",
		Status:   http.StatusBadRequest,
		Detail:   e.Error(),
		Instance: c.Path(),
		Method:   c.Method(),
		Datetime: time.Now().UTC(),
	}

	switch err := e.(type) {
	case ship.HTTPServerError:
		pd.Status = err.Code
	case *ship.HTTPServerError:
		pd.Status = err.Code
	// case *validate.TranError:
	//	pd.Title = "参数校验错误"
	case *time.ParseError:
		pd.Detail = "时间格式错误，正确格式：" + err.Layout
	case *net.ParseError:
		pd.Detail = err.Text + " 不是有效的 " + err.Type
	case base64.CorruptInputError:
		pd.Detail = "base64 编码错误：" + err.Error()
	case *json.SyntaxError:
		pd.Detail = "请求报错必须是 JSON 格式"
	case *json.UnmarshalTypeError:
		pd.Detail = err.Field + " 收到无效的数据类型"
	case *strconv.NumError:
		var msg string
		if sn := strings.SplitN(err.Func, "Parse", 2); len(sn) == 2 {
			msg = err.Num + " 不是 " + strings.ToLower(sn[1]) + " 类型"
		} else {
			msg = "类型错误：" + err.Num
		}
		pd.Detail = msg
	case *http.MaxBytesError:
		limit := strconv.FormatInt(err.Limit, 10)
		pd.Detail = "请求报文超过 " + limit + " 个字节限制"
		pd.Status = http.StatusRequestEntityTooLarge
	// case mongo.CommandError:
	//	pd.Status = http.StatusInternalServerError
	//	pd.Title = "内部错误"
	//	pd.Detail = "内部错误请联系管理员"
	// case topology.ServerSelectionError:
	//	pd.Status = http.StatusServiceUnavailable
	//	pd.Title = "内部错误"
	//	pd.Detail = "内部错误请联系管理员"
	//	if errors.Is(err.Wrapped, context.DeadlineExceeded) {
	//		pd.Detail = "内部服务超时"
	//	}
	//	c.Errorf("数据库连接错误", slog.Any("error", e))
	default:
		switch {
		// case errors.Is(err, primitive.ErrInvalidHex):
		//	pd.Detail = "不是一个有效的 ID"
		// case errors.Is(err, mongo.ErrNoDocuments):
		//	pd.Detail = "数据不存在"
		case errors.Is(err, ship.ErrSessionNotExist), errors.Is(err, ship.ErrInvalidSession):
			pd.Status = http.StatusUnauthorized
			pd.Detail = "认证无效"
		}
	}

	_ = c.JSON(pd.Status, pd)
}
