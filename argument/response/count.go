package response

import "gorm.io/gen/field"

type NameCount struct {
	Name  string `json:"name"  gorm:"column:name"`
	Count int64  `json:"count" gorm:"column:count"`
}

type NameCounts []*NameCount

func (nc NameCounts) Aliases() (string, string, field.OrderExpr) {
	const count = "count"
	expr := field.NewField("", count)
	return "name", "count", expr // 与 gorm tag 保持一致
}
