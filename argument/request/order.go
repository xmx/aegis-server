package request

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/gormcond"
)

type PageKeywordOrder struct {
	PageKeyword
	Order Orders `json:"order" form:"order" query:"order" validate:"lte=100"`
}

type Orders []*Order

func (ods Orders) Orders() []*gormcond.ColumnDesc {
	size := len(ods)
	ret := make([]*gormcond.ColumnDesc, 0, size)
	for _, od := range ods {
		if col := od.Column; col != "" {
			ret = append(ret, &gormcond.ColumnDesc{Column: col, Desc: od.Desc})
		}
	}
	return ret
}

type Order struct {
	Column string `json:"column"`
	Desc   bool   `json:"desc"`
}

// UnmarshalBind 从 query 参数中解析排序规则。
//
// 支持常规的 JSON 格式：
//   - {"column":"age"}                 ORDER BY age
//   - {"column":"age","desc":false}    ORDER BY age
//   - {"column":"age","desc":true}     ORDER BY age DESC
//
// 支持 <field>:<desc> 格式（因为 JSON 格式放在 query 中太冗长）：
//
//	<desc> 是未定义的值，或者是空（此时冒号可以省略），均按照升序排序：
//		- age:hello    ORDER BY age
//		- age:balabala ORDER BY age
//		- age:         ORDER BY age
//		- age	       ORDER BY age
//	<desc> 可以是枚举：asc、desc：
//		- age:asc      ORDER BY age
//		- age:desc     ORDER BY age DESC
//	<desc> 可以是数字，小于 0 时代表降序，否则代表升序：
//		- age:0        ORDER BY age
//		- age:123      ORDER BY age
//		- age:-1       ORDER BY age DESC
//		- age:-456     ORDER BY age DESC
//
// 错误格式：
//   - :desc
//   - :巴拉巴拉
func (o *Order) UnmarshalBind(str string) error {
	if data := []byte(str); json.Valid(data) {
		return json.Unmarshal(data, o)
	}

	idx := strings.LastIndex(str, ":")
	if idx < 0 {
		o.Column = str
		return nil
	} else if idx == 0 {
		return errcode.ErrRequiredColumn
	}

	o.Column = str[:idx]
	value := str[idx+1:]
	if sorted := strings.ToLower(str[idx+1:]); sorted == "asc" || sorted == "desc" {
		o.Desc = sorted == "desc"
		return nil
	}

	num, _ := strconv.ParseInt(value, 10, 64)
	o.Desc = num < 0

	return nil
}
