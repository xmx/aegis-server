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

//  =          eq       1 string int bool time
// !=          ne       1 string int bool time
//  >          gt       1 string int      time
//  <          lt       1 string int      time
// >=          gte      1 string int      time
// <=          lte      1 string int      time
// IN          in       n string int      time
// NOT IN      notin    n string int      time
// LIKE        like     1 string
// NOT LIKE    notlike  1 string
// BETWEEN     btw      2 string int      time
// NOT BETWEEN notbtw   2 string int      time
// REGEX       regex    1 string
// NOT REGEX   notregex 1 string
