package request

type PageKeyword struct {
	Page
	OptionalKeyword
}

type OptionalKeyword struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"lte=100"`
}

func (o OptionalKeyword) Like() string {
	if kw := o.String(); kw != "" {
		return "%" + kw + "%"
	}
	return ""
}

func (o OptionalKeyword) LLike() string {
	if kw := o.String(); kw != "" {
		return "%" + kw
	}
	return ""
}

func (o OptionalKeyword) RLike() string {
	if kw := o.String(); kw != "" {
		return kw + "%"
	}
	return ""
}

func (o OptionalKeyword) String() string {
	return o.Keyword
}
