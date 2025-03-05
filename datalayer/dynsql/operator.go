package dynsql

type Operators []Operator

type Operator interface {
	// NArgs 该操作至少需要几个参数。
	// <=0 标表示不需要参数
	NArgs() int
	Info() (id, name string)
	String() string
}

type operator struct {
	id   string
	name string
	args int
}

func (o operator) NArgs() int {
	return o.args
}

func (o operator) Info() (string, string) {
	return o.id, o.name
}

func (o operator) String() string {
	return o.id
}

var (
	Eq         = operator{id: "eq", name: "=", args: 1}
	Neq        = operator{id: "neq", name: "≠", args: 1}
	Gt         = operator{id: "gt", name: ">", args: 1}
	Gte        = operator{id: "gte", name: "≥", args: 1}
	Lt         = operator{id: "lt", name: "<", args: 1}
	Lte        = operator{id: "lte", name: "≤", args: 1}
	Like       = operator{id: "like", name: "LIKE", args: 1}
	NotLike    = operator{id: "notlike", name: "NOT LIKE", args: 1}
	Between    = operator{id: "between", name: "BETWEEN", args: 2}
	NotBetween = operator{id: "notbetween", name: "NOT BETWEEN", args: 2}
	In         = operator{id: "in", name: "IN", args: 1}
	NotIn      = operator{id: "notin", name: "NOT IN", args: 1}
	Null       = operator{id: "null", name: "NULL", args: 0}

	operatorMaps = map[string]Operator{
		Eq.id:         Eq,
		Neq.id:        Neq,
		Gt.id:         Gt,
		Gte.id:        Gte,
		Lt.id:         Lt,
		Lte.id:        Lte,
		Like.id:       Like,
		NotLike.id:    NotLike,
		Between.id:    Between,
		NotBetween.id: NotBetween,
		In.id:         In,
		NotIn.id:      NotIn,
		Null.id:       Null,
	}
)

func Lookup(id string) Operator {
	return operatorMaps[id]
}
