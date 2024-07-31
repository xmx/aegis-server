package jslib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/dop251/goja"
)

type formatter interface {
	format(call goja.FunctionCall) ([]byte, error)
}

type normalFormat struct{}

func (n normalFormat) format(call goja.FunctionCall) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, val := range call.Arguments {
		if err := n.parse(buf, val); err != nil {
			return nil, err
		}
		buf.WriteByte(' ')
	}

	return buf.Bytes(), nil
}

func (normalFormat) parse(buf *bytes.Buffer, val goja.Value) error {
	switch {
	case goja.IsUndefined(val), goja.IsNull(val):
		buf.WriteString(val.String())
		return nil
	}

	export := val.Export()
	switch v := export.(type) {
	case fmt.Stringer:
		buf.WriteString(v.String())
	case string:
		buf.WriteString(v)
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
	case float64:
		buf.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
	case bool:
		buf.WriteString(strconv.FormatBool(v))
	case []byte:
		str := base64.StdEncoding.EncodeToString(v)
		buf.WriteString(str)
	case func(goja.FunctionCall) goja.Value:
		buf.WriteString("<Function>")
	case goja.ArrayBuffer:
		bs := v.Bytes()
		str := base64.StdEncoding.EncodeToString(bs)
		buf.WriteString(str)
	default:
		tmp := new(bytes.Buffer)
		if err := json.NewEncoder(tmp).Encode(val); err != nil || tmp.Len() == 0 {
			return err
		}
		// 标准库 JSON 序列化后会在最后面加换行符，这样输出到前端会有个换行，
		// 下面的操作就是去除 JSON 字符串后面的换行符。当然不处理也无伤大雅。
		_, _ = buf.ReadFrom(io.LimitReader(tmp, int64(tmp.Len()-1)))
	}

	return nil
}
