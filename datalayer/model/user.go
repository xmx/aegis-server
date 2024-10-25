package model

import "time"

type User struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Disabled  bool      `json:"disabled"   gorm:"column:disabled;comment:是否禁用"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:now(3);comment:更新时间"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:now(3);comment:创建时间"`
}

func (User) TableName() string {
	return "user"
}
