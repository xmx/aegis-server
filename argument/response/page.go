package response

type Page[T any] struct {
	Page    int64 `json:"page"`
	Size    int64 `json:"size"`
	Count   int64 `json:"count"`
	Records []T   `json:"records"`
}
