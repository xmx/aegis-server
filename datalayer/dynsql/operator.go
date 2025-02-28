package dynsql

type Operators []Operator

type Operator interface {
	OpInfo() (string, string)
}

type operator struct {
	id   string
	name string
}

func (o operator) OpInfo() (string, string) {
	return o.id, o.name
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
	NotNull    = operator{id: "notnull", name: "NOT NULL"}

	opMaps = map[string]Operator{
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
		NotNull.id:    NotNull,
	}
)

func Lookup(id string) Operator {
	return opMaps[id]
}
