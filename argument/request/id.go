package request

import "go.mongodb.org/mongo-driver/v2/bson"

type ObjectID struct {
	ID string `json:"id" query:"id" form:"id" validate:"required,mongodb"`
}

func (o ObjectID) OID() bson.ObjectID {
	id, _ := bson.ObjectIDFromHex(o.ID)
	return id
}

type ObjectIDs struct {
	ID []string `json:"id" query:"id" form:"id" validate:"gte=1,lte=1000,dive,required,mongodb"`
}

func (o ObjectIDs) OIDs() []bson.ObjectID {
	ids := make([]bson.ObjectID, 0, len(o.ID))

	for _, str := range o.ID {
		if id, err := bson.ObjectIDFromHex(str); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}
