package middle

import (
	"errors"
	"fmt"
	"time"

	"github.com/xgfone/ship/v5"
)

type RecordOptions struct {
	// Description 简短的描述路由功能（必填）。
	// 例如：添加证书、删除 agent 节点 等等。
	Description string

	// Latency 如果该值大于 0 则代表当函数处理延迟大于等于 Latency 时记录日志，否则不记录。
	Latency time.Duration

	// RecordRequestBytes 中间件记录请求报文大小，小于等于 0 代表不记录。
	RecordRequestBytes int

	// RecordResponseBytes 中间件记录响应报文大小，小于等于 0 代表不记录。
	RecordResponseBytes int
}

func NewRecord(desc string) RecordOptions {
	return RecordOptions{
		Description: desc,
	}
}

type RouteInfo struct {
	RecordOptions RecordOptions
}

func CheckRouteInfo(routes []ship.Route) error {
	var errs []error
	for _, route := range routes {
		mth, uri := route.Method, route.Path
		info, ok := route.Data.(*RouteInfo)
		if !ok || info == nil {
			e := fmt.Errorf(`%s %s: 缺少定义中间件信息，例如：r.Route(%s).Data(此处请填写中间件信息).%s(handle)`, mth, uri, uri, mth)
			errs = append(errs, e)
			continue
		}
		if info.RecordOptions.Description == "" {
			e := fmt.Errorf("%s %s: 缺少路由说明", mth, uri)
			errs = append(errs, e)
		}
	}

	return errors.Join(errs...)
}
