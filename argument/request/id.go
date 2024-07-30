package request

import "strconv"

type Int64ID struct {
	ID int64 `json:"id,string" form:"id" query:"id" validate:"required"`
}

type Int64IDs struct {
	ID []int64 `json:"id" form:"id" query:"id" validate:"gte=1,lte=1000,dive,required"`
}

type OptionalInt64IDs struct {
	ID []int64 `json:"id" form:"id" query:"id" validate:"lte=1000,dive,required"`
}

type OptionalInt64ID struct {
	ID NullInt64 `json:"id" query:"id"`
}

type NullInt64 struct {
	value int64
	valid bool
}

func (n *NullInt64) Get() (int64, bool) {
	return n.value, n.valid
}

func (n *NullInt64) UnmarshalText(str []byte) error {
	return n.UnmarshalBind(string(str))
}

func (n *NullInt64) UnmarshalBind(str string) error {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	n.value = num
	n.valid = true

	return nil
}
