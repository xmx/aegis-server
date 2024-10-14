package condition

var (
	Eq         = Operator{name: "eq"}
	Neq        = Operator{name: "neq"}
	Gt         = Operator{name: "gt"}
	Gte        = Operator{name: "gte"}
	Lt         = Operator{name: "lt"}
	Lte        = Operator{name: "lte"}
	Like       = Operator{name: "like"}
	NotLike    = Operator{name: "not-like"}
	Regex      = Operator{name: "regex"}
	NotRegex   = Operator{name: "not-regex"}
	Between    = Operator{name: "between"}
	NotBetween = Operator{name: "not-between"}
	In         = Operator{name: "in"}
	NotIn      = Operator{name: "not-in"}
)

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
