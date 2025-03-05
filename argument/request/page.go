package request

type Pages struct {
	Page int64 `json:"page" query:"page" validate:"gte=0"`
	Size int64 `json:"size" query:"size" validate:"lte=1000"`
}

func (p Pages) PageSize() (int64, int64) {
	return p.Page, p.Size
}
