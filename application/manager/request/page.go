package request

type Pages struct {
	Page int64  `json:"page" form:"page" query:"page" validate:"gte=0"`
	Size int64  `json:"size" form:"size" query:"size" validate:"gte=0,lte=1000"`
	Q    string `json:"q"    form:"q"    query:"q"    validate:"lte=1000"`
}

type IDPages struct {
	ObjectID
	Pages
}
