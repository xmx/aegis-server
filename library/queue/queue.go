package queue

type Queuer[T any] interface {
	Enqueue(T)
	Dequeue() (T, bool)
}
