package model

import (
	"sort"
	"time"
)

type Menu struct {
	ID        int64     `json:"id,string"        gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	ParentID  int64     `json:"parent_id,string" gorm:"column:parent_id;type:bigint;comment:父ID"`
	Name      string    `json:"name"             gorm:"column:name;type:varchar(20);index:;comment:菜单名"`
	Key       string    `json:"key"              gorm:"column:key;type:varchar(100);index:;comment:菜单标识"`
	Data      string    `json:"data"             gorm:"column:data;type:text;comment:自定义数据"`
	Hide      bool      `json:"hide"             gorm:"column:hide;comment:是否隐藏"`
	Icon      string    `json:"icon"             gorm:"column:icon;type:text;comment:图标"`
	Folder    bool      `json:"folder"           gorm:"column:folder;comment:是否目录"`
	Order     int64     `json:"order"            gorm:"column:order;type:bigint;comment:次序"`
	UpdatedAt time.Time `json:"updated_at"       gorm:"column:updated_at;not null;default:now(3);comment:更新时间"`
	CreatedAt time.Time `json:"created_at"       gorm:"column:created_at;not null;default:now(3);comment:创建时间"`
}

func (Menu) TableName() string {
	return "menu"
}

func (m Menu) Node() *MenuNode {
	return &MenuNode{Menu: m}
}

func (m Menu) IsRoot() bool {
	return m.ParentID == 0
}

type MenuNode struct {
	Menu
	Children MenuNodes `json:"children,omitempty"`
}

type MenuNodes []*MenuNode

func (mns MenuNodes) Len() int {
	return len(mns)
}

func (mns MenuNodes) Less(i, j int) bool {
	mi, mj := mns[i], mns[j]
	if mi.Folder != mj.Folder {
		return mi.Folder
	}
	if mi.Order != mj.Order {
		return mi.Order < mj.Order
	}
	return mi.Name < mj.Name
}

func (mns MenuNodes) Swap(i, j int) {
	mns[i], mns[j] = mns[j], mns[i]
}

func (mns MenuNodes) Sort() {
	mns.sort(mns)
}

func (mns MenuNodes) sort(nodes MenuNodes) {
	sort.Sort(nodes)
	for _, mn := range nodes {
		mns.sort(mn.Children)
	}
}
