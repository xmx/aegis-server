package request

import "go.mongodb.org/mongo-driver/v2/bson"

type ObjectID struct {
	ID string `json:"id" form:"id" query:"id" validate:"mongodb"`
}

func (o ObjectID) OID() bson.ObjectID {
	id, _ := bson.ObjectIDFromHex(o.ID)
	return id
}

type ObjectIDs struct {
	ID []string `json:"id" form:"id" query:"id" validate:"dive,mongodb"`
}

func (o ObjectIDs) OIDs() []bson.ObjectID {
	ids := make([]bson.ObjectID, 0, len(o.ID))
	for _, s := range o.ID {
		if id, err := bson.ObjectIDFromHex(s); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}
