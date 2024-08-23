package model

import "time"

type ConfigServer struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement"` // 数据库 ID。
	Enabled   bool      `json:"enabled"    gorm:"column:enabled"`                     // 是否启用该配置。
	Addr      string    `json:"addr"       gorm:"column:addr;type:varchar(100)"`      // 服务监听地址。
	Static    string    `json:"static"     gorm:"column:static;type:varchar(255)"`    // 静态资源目录。
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:now(3)"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:now(3)"`
}

func (ConfigServer) TableName() string {
	return "config_server"
}
