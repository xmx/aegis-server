package request

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/datalayer/dynsql"
)

type Cond struct {
	Wheres CondWheres `json:"wheres" query:"wheres" form:"wheres" validate:"lte=100"`
	Orders CondOrders `json:"orders" query:"orders" form:"orders" validate:"lte=10"`
}

type CondWheres []*CondWhere

type CondWhere struct {
	ID     string   `json:"id"     validate:"lte=100"`
	Op     string   `json:"op"     validate:"oneof=eq neq gt gte lt lte like notlike between notbetween in notin null"`
	Values []string `json:"values" validate:"lte=100,dive,lte=100"`
}

// UnmarshalBind 从 query 参数中解析查询条件。
//
// 支持常规的 JSON 格式：
//   - {"id":"age","op":"eq","values":["16"]}   WHERE age = 16
//   - {"id":"age","op":"in","values":["16", "17"]}   WHERE age IN (16, 17)
//   - {"id":"age","op":"null","values":[]}   WHERE age IS NULL
//   - {"id":"age","op":"null","values":["false"]}   WHERE age IS NOT NULL
//
// 支持 <id>:<desc> 格式（因为 JSON 格式放在 query 中太冗长）：
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
func (cw *CondWhere) UnmarshalBind(str string) error {
	if data := []byte(str); json.Valid(data) {
		return json.Unmarshal(data, cw)
	}

	// created_at:gte:2024-10-14T06:31:20.620Z
	// size:between:1024,4096
	// hobbies:in:sing,dance,rap

	sn := strings.SplitN(str, ":", 3)
	if len(sn) != 3 {
		return errcode.ErrRequiredColumn
	}
	name, opera, val := sn[0], sn[1], sn[2]
	cw.ID, cw.Op = name, opera
	op := dynsql.Lookup(opera)
	args := op.NArgs()
	if args <= 0 {
		cw.Values = []string{val}
	}

	vals := strings.Split(val, ",")
	if len(vals) < args {
		return errcode.ErrRequiredArgs
	}
	cw.Values = vals

	return nil
}

type CondOrders []*CondOrder

type CondOrder struct {
	ID   string `json:"id"   validate:"lte=100"`
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
func (co *CondOrder) UnmarshalBind(str string) error {
	if data := []byte(str); json.Valid(data) {
		return json.Unmarshal(data, co)
	}

	id, val, _ := strings.Cut(str, ":")
	var desc bool
	if val != "" {
		val = strings.ToLower(val)
		if val == "asc" || val == "desc" {
			desc = val == "desc"
		} else {
			num, _ := strconv.ParseInt(val, 10, 64)
			desc = num < 0
		}
	}
	co.ID, co.Desc = id, desc

	return nil
}
