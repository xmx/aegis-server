package model

import "time"

type ConfigServer struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Enabled   bool      `json:"enabled"    gorm:"column:enabled;comment:是否启用"`
	Addr      string    `json:"addr"       gorm:"column:addr;type:varchar(100);comment:监听地址"`
	Static    string    `json:"static"     gorm:"column:static;type:varchar(255);comment:静态资源目录"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:now(3);comment:更新时间"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:now(3);comment:创建时间"`
}

func (ConfigServer) TableName() string {
	return "config_server"
}
