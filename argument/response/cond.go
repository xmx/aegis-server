package response

import "github.com/xmx/aegis-server/argument/bizdata"

type Cond struct {
	Orders []*bizdata.FieldComment `json:"orders"`
}
