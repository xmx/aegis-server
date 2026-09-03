package request

import "go.mongodb.org/mongo-driver/v2/bson"

type ObjectID struct {
	ID TextID `json:"id" form:"id" query:"id" param:"id" validate:"required,mongodb"`
}

type TextID string

func (tid TextID) Raw() string {
	return string(tid)
}

func (tid TextID) ObjectID() (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(tid.Raw())
}

func (tid TextID) UnsafeID() bson.ObjectID {
	id, _ := tid.ObjectID()
	return id
}

type PID struct {
	PID int32 `json:"pid" form:"pid" query:"pid" validate:"gte=0"`
}
