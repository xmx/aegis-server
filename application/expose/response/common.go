package response

type FieldValues[F, E any] struct {
	Field  F   `json:"field"  bson:"field"`
	Values []E `json:"values" bson:"values"`
}
