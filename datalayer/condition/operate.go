package condition

var (
	Eq         = NewOperator("eq")
	Neq        = NewOperator("neq")
	Gt         = NewOperator("gt")
	Gte        = NewOperator("gte")
	Lt         = NewOperator("lt")
	Lte        = NewOperator("lte")
	Like       = NewOperator("like")
	NotLike    = NewOperator("not_like")
	Regex      = NewOperator("regex")
	NotRegex   = NewOperator("not_regex")
	Between    = NewOperator("between")
	NotBetween = NewOperator("not_between")
	In         = NewOperator("in")
	NotIn      = NewOperator("not_in")
)

func NewOperator(name string) Operator {
	return Operator{name: name}
}

type Operator struct {
	name string
}

func (o Operator) String() string {
	return o.name
}

func (o Operator) IsNoop() bool {
	return o.name == ""
}

type operators []Operator

func (ops operators) NameMap() map[string]Operator {
	hm := make(map[string]Operator, 8)
	for _, op := range ops {
		hm[op.name] = op
	}
	return hm
}
