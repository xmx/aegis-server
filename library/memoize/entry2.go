package memoize

type entry2[V, E any] struct {
	v V
	e E
}

func (e *entry2[V, E]) load() (V, E) {
	return e.v, e.e
}
