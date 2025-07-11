package model

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Oplog struct {
	ID          bson.ObjectID  `json:"id"                 bson:"_id,omitempty"`      // 日志 ID
	RouteName   string         `json:"route_name"         bson:"route_name"`         // 业务名
	Operator    *Operator      `json:"operator,omitempty" bson:"operator,omitempty"` // 操作用户
	Request     *OplogRequest  `json:"request,omitempty"  bson:"request,omitempty"`  // 请求信息
	Response    *OplogResponse `json:"response,omitempty" bson:"response,omitempty"` // 响应相关信息
	Succeed     bool           `json:"succeed"            bson:"succeed"`            // 业务处理是否成功
	Reason      string         `json:"reason,omitempty"   bson:"reason,omitempty"`   // 如果出错，错误原因
	RequestedAt time.Time      `json:"requested_at"       bson:"requested_at"`       // 前端请求时间
	FinishedAt  time.Time      `json:"finished_at"        bson:"finished_at"`        // 后端处理结束时间
	Elapsed     time.Duration  `json:"elapsed"            bson:"elapsed"`            // 处理耗时
}

type OplogRequest struct {
	Method     string      `json:"method"      bson:"method"`      // 方法
	Path       string      `json:"path"        bson:"path"`        // 请求路径
	Query      string      `json:"query"       bson:"query"`       // Query 参数
	RemoteAddr string      `json:"remote_addr" bson:"remote_addr"` // 客户端地址，从 Header 中获取
	DirectAddr string      `json:"direct_addr" bson:"direct_addr"` // 直连地址，可能是反向代理的地址
	Header     http.Header `json:"header"      bson:"header"`      // Header
	Body       []byte      `json:"body"        bson:"body"`        // 请求报文（超过4096的部分截断丢弃）
}

type OplogResponse struct {
	StatusCode int         `json:"status_code" bson:"status_code"`
	Header     http.Header `json:"header"      bson:"header"`
	Body       []byte      `json:"body"        bson:"body"`
}
