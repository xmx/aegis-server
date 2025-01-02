package response

import "gorm.io/gen/field"

type NameCount struct {
	Name  string `json:"name"  gorm:"column:name"`
	Count int64  `json:"count" gorm:"column:count"`
}

type NameCounts []*NameCount

func (nc NameCounts) Aliases() (name, count field.Field) {
	// 要与 NameCount gorm tag 保持一致
	name = field.NewField("", "name")
	count = field.NewField("", "count")
	return
}
