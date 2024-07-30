package request

type Int64ID struct {
	ID int64 `json:"id,string" form:"id" query:"id" validate:"required"`
}

type Int64IDs struct {
	ID []int64 `json:"id" form:"id" query:"id" validate:"gte=1,lte=1000,dive,required"`
}
