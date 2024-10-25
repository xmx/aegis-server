package model

import "time"

type Role struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Name      string    `json:"name"       gorm:"column:name;type:varchar(20);uniqueIndex;comment:角色名"`
	Admin     bool      `json:"admin"      gorm:"column:admin;comment:是否超管"`
	Disabled  bool      `json:"disabled"   gorm:"column:disabled;comment:是否禁用"`
	Remark    string    `json:"remark"     gorm:"column:remark;type:text;comment:备注说明"`
	CreatedAt time.Time `json:"created_at" gorm:"column:updated_at;not null;default:now(3);comment:更新时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:created_at;not null;default:now(3);comment:创建时间"`
}

func (Role) TableName() string {
	return "role"
}

func (r Role) EnabledAdministrator() bool {
	return !r.Disabled && r.Admin
}
