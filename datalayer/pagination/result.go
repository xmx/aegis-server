package pagination

type Result[E any] struct {
	Page    int64 `json:"page"`
	Size    int64 `json:"size"`
	Count   int64 `json:"count"`
	Records []E   `json:"records"`
}
