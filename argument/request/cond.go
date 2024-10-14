package request

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/datalayer/condition"
)

type PageKeywordCond struct {
	PageKeyword
	Order CondOrders `json:"order" form:"order" query:"order" validate:"lte=100"`
	Where CondWheres `json:"where" form:"where" query:"where" validate:"lte=100"`
}

type CondOrders []*CondOrder

func (ods CondOrders) Orders() condition.OrderInputs {
	size := len(ods)
	ret := make(condition.OrderInputs, 0, size)
	for _, od := range ods {
		if col := od.Name; col != "" {
			data := &condition.OrderInput{Name: col, Desc: od.Desc}
			ret = append(ret, data)
		}
	}
	return ret
}

type CondOrder struct {
	Name string `json:"name"`
	Desc bool   `json:"desc"`
}

// UnmarshalBind 从 query 参数中解析排序规则。
//
// 支持常规的 JSON 格式：
//   - {"name":"age"}                 ORDER BY age
//   - {"name":"age","desc":false}    ORDER BY age
//   - {"name":"age","desc":true}     ORDER BY age DESC
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
func (o *CondOrder) UnmarshalBind(str string) error {
	if data := []byte(str); json.Valid(data) {
		return json.Unmarshal(data, o)
	}

	idx := strings.LastIndex(str, ":")
	if idx < 0 {
		o.Name = str
		return nil
	} else if idx == 0 {
		return errcode.ErrRequiredColumn
	}

	o.Name = str[:idx]
	value := str[idx+1:]
	if sorted := strings.ToLower(str[idx+1:]); sorted == "asc" || sorted == "desc" {
		o.Desc = sorted == "desc"
		return nil
	}

	num, _ := strconv.ParseInt(value, 10, 64)
	o.Desc = num < 0

	return nil
}

type CondWhere struct {
	Name    string   `json:"name"    validate:"required"`
	Operate string   `json:"operate" validate:"required"`
	Values  []string `json:"values"`
}

func (c *CondWhere) UnmarshalBind(param string) error {
	return json.Unmarshal([]byte(param), c)
}

type CondWheres []*CondWhere

func (cws CondWheres) Wheres() condition.WhereInputs {
	size := len(cws)
	ret := make(condition.WhereInputs, 0, size)
	for _, cw := range cws {
		if col := cw.Name; col != "" {
			data := &condition.WhereInput{
				Name:    col,
				Operate: cws.operator(cw.Operate),
				Values:  cw.Values,
			}
			ret = append(ret, data)
		}
	}
	return ret
}

func (CondWheres) operator(s string) condition.Operator {
	hm := map[string]condition.Operator{
		"=":   condition.Eq,
		"!=":  condition.Neq,
		">":   condition.Gt,
		">=":  condition.Gte,
		"<":   condition.Lt,
		"<=":  condition.Lte,
		"~=":  condition.Like,
		"!~=": condition.NotLike,
		"$=":  condition.Regex,
		"!$=": condition.NotRegex,
		"*=":  condition.In,
		"!*=": condition.NotIn,
		"^=":  condition.Between,
		"!^=": condition.NotBetween,
	}
	return hm[s]
}
