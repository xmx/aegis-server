package condition

var (
	Eq         = NewOperator("eq")
	Neq        = NewOperator("neq")
	Gt         = NewOperator("gt")
	Gte        = NewOperator("gte")
	Lt         = NewOperator("lt")
	Lte        = NewOperator("lte")
	Like       = NewOperator("like")
	NotLike    = NewOperator("not-like")
	Regex      = NewOperator("regex")
	NotRegex   = NewOperator("not-regex")
	Between    = NewOperator("between")
	NotBetween = NewOperator("not-between")
	In         = NewOperator("in")
	NotIn      = NewOperator("not-in")
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
