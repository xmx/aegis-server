package response

import "github.com/xmx/aegis-server/argument/gormcond"

type Cond struct {
	Wheres []any                     `json:"wheres,omitempty"`
	Orders []*gormcond.ColumnComment `json:"orders,omitempty"`
}
