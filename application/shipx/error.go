package shipx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/problem"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/topology"
)

func NotFound(*ship.Context) error { return ship.ErrNotFound }

func HandleError(c *ship.Context, e error) {
	he := HTTPError(e)
	if he == nil {
		return
	}

	status := he.Code
	if status < http.StatusBadRequest || status >= 1000 {
		status = http.StatusBadRequest
	}
	r := c.Request()

	body := &problem.Details{
		Status:   status,
		Title:    http.StatusText(status),
		Detail:   he.Error(),
		Instance: r.URL.Path,
		Method:   r.Method,
		Host:     r.Host,
	}

	_ = c.JSON(status, body)
}

func HTTPError(err error) *ship.HTTPServerError {
	switch e := err.(type) {
	case nil:
		return nil
	case ship.HTTPServerError:
		return new(e)
	case *ship.HTTPServerError:
		return e
	case *http.MaxBytesError:
		return new(ship.NewHTTPServerError(http.StatusRequestEntityTooLarge, "请求正文过大"))
	case *time.ParseError:
		return new(ship.ErrBadRequest.Newf("时间格式错误"))
	case *net.ParseError:
		return new(ship.ErrBadRequest.Newf("无效的 %s：%s", e.Type, e.Text))
	case base64.CorruptInputError:
		return new(ship.ErrBadRequest.Newf("第 %d 字节处存在非法的 Base64 数据", e))
	case *json.SyntaxError:
		return new(ship.ErrBadRequest.Newf("第 %d 字节处存在 JSON 语法错误：%s", e.Offset, e))
	case *json.UnmarshalTypeError:
		if e.Struct != "" || e.Field != "" {
			return new(ship.ErrBadRequest.Newf("[JSON 解析] 字段 %s 需要%s类型，但收到: %s", e.Field, e.Type.String(), e.Value))
		}
		return new(ship.ErrBadRequest.Newf("[JSON 解析] 无法将 %s 转换为 %s 类型", e.Value, e.Type.String()))
	case *strconv.NumError:
		return new(ship.ErrBadRequest.Newf("无法将 %q 转换为数字（原因：%v）", e.Num, e.Err))
	case mongo.CommandError:
		if e.HasErrorCode(11000) {
			return new(ship.ErrBadRequest.Newf("数据已存在"))
		}
		return new(ship.ErrBadRequest.Newf("[数据库] 操作错误：%s", e.Message))
	case topology.ServerSelectionError:
		if e.Wrapped != nil {
			return new(ship.ErrInternalServerError.Newf("[数据库] 服务器选择错误：%s，当前拓扑结构：%s", e.Wrapped, e.Desc))
		}
		return new(ship.ErrInternalServerError.Newf("[数据库] 服务器选择错误，当前拓扑结构：%s", e.Desc))
	case mongo.WriteException:
		if e.HasErrorCode(11000) {
			return new(ship.ErrBadRequest.Newf("数据已存在"))
		}
		return new(ship.ErrBadRequest.New(e))
	}

	switch {
	case errors.Is(err, bson.ErrInvalidHex):
		return new(ship.ErrBadRequest.Newf("字符串不是有效的 ObjectID"))
	case errors.Is(err, mongo.ErrNoDocuments):
		return new(ship.ErrBadRequest.Newf("数据不存在"))
	case errors.Is(err, ship.ErrSessionNotExist), errors.Is(err, ship.ErrInvalidSession):
		return new(ship.ErrUnauthorized.Newf("认证无效"))
	case errors.Is(err, context.DeadlineExceeded):
		return new(ship.ErrBadRequest.Newf("数据超时"))
	case errors.Is(err, context.Canceled):
		return new(ship.ErrBadRequest.Newf("数据取消"))
	}

	return new(ship.ErrBadRequest.New(err))
}
