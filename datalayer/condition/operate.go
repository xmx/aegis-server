package condition

var (
	Eq         = condOp{name: "eq"}
	Neq        = condOp{name: "neq"}
	Gt         = condOp{name: "gt"}
	Gte        = condOp{name: "gte"}
	Lt         = condOp{name: "lt"}
	Lte        = condOp{name: "lte"}
	Like       = condOp{name: "like"}
	NotLike    = condOp{name: "not-like"}
	Regex      = condOp{name: "regex"}
	NotRegex   = condOp{name: "not-regex"}
	Between    = condOp{name: "between"}
	NotBetween = condOp{name: "not-between"}
	In         = condOp{name: "in"}
	NotIn      = condOp{name: "not-in"}
)

type condOp struct {
	name string
}

func (o condOp) isNoop() bool {
	return o.name == ""
}
