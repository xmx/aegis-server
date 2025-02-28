package dynsql

type Operators []Operator

type Operator interface {
	OpInfo() (id, name string)
	String() string
}

type operator struct {
	id   string
	name string
}

func (o operator) OpInfo() (string, string) {
	return o.id, o.name
}

func (o operator) String() string {
	return o.id
}

var (
	Eq         = operator{id: "eq", name: "="}
	Neq        = operator{id: "neq", name: "≠"}
	Gt         = operator{id: "gt", name: ">"}
	Gte        = operator{id: "gte", name: "≥"}
	Lt         = operator{id: "lt", name: "<"}
	Lte        = operator{id: "lte", name: "≤"}
	Like       = operator{id: "like", name: "LIKE"}
	NotLike    = operator{id: "notlike", name: "NOT LIKE"}
	Between    = operator{id: "between", name: "BETWEEN"}
	NotBetween = operator{id: "notbetween", name: "NOT BETWEEN"}
	In         = operator{id: "in", name: "IN"}
	NotIn      = operator{id: "notin", name: "NOT IN"}
	Null       = operator{id: "null", name: "NULL"}

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
