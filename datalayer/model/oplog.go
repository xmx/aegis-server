package model

import (
	"net/http"
	"net/url"
	"time"
)

type Oplog struct {
	ID         int64       `json:"id,string,omitempty" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Name       string      `json:"name"                gorm:"column:name;type:varchar(50);comment:业务名"`
	Host       string      `json:"host"                gorm:"column:host;type:varchar(50);comment:主机"`
	Method     string      `json:"method"              gorm:"column:method;type:varchar(10);comment:方法"`
	Path       string      `json:"path"                gorm:"column:path;type:varchar(255);comment:路径"`
	Query      url.Values  `json:"query,omitempty"     gorm:"column:query;type:json;serializer:json;comment:查询参数"`
	Body       []byte      `json:"body,omitempty"      gorm:"column:body;type:blob;comment:报文"`
	Header     http.Header `json:"header,omitempty"    gorm:"column:header;type:json;serializer:json;comment:Header"`
	ClientIP   string      `json:"client_ip"           gorm:"column:client_ip;type:varchar(50);comment:客户端IP"`
	DirectIP   string      `json:"direct_ip"           gorm:"column:direct_ip;type:varchar(50);comment:直接客户端IP"`
	Succeed    bool        `json:"succeed"             gorm:"column:succeed;comment:是否成功"`
	Reason     string      `json:"reason,omitempty"    gorm:"column:reason;type:text;comment:原因"`
	AccessedAt time.Time   `json:"accessed_at"         gorm:"column:accessed_at;not null;default:now(3);index:,sort:desc;comment:访问起始时间"`
	FinishedAt time.Time   `json:"finished_at"         gorm:"column:finished_at;not null;default:now(3);comment:访问结束时间"`
}

func (Oplog) TableName() string {
	return "oplog"
}
