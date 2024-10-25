package model

type RoleMenu struct {
	ID     int64 `json:"-"              gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	RoleID int64 `json:"role_id,string" gorm:"column:role_id;index:idx_role_menu_id;comment:角色ID"`
	MenuID int64 `json:"menu_id,string" gorm:"column:menu_id;index:idx_role_menu_id;comment:菜单ID"`
}

func (RoleMenu) TableName() string {
	return "role_menu"
}
