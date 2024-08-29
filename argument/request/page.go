package request

type Page struct {
	Page int64 `json:"page" query:"page" validate:"gte=0"`
	Size int64 `json:"size" query:"size" validate:"lte=1000"`
}
