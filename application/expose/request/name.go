package request

type Names struct {
	Name string `json:"name" form:"name" query:"name" validate:"required"`
}
