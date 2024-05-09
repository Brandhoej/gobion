package structures

type Queue[T any] []T

func (queue *Queue[T]) IsEmpty() bool {
	return len(*queue) == 0
}

func (queue *Queue[T]) Enqueue(elements ...T) {
	*queue = append(*queue, elements...)
}

func (queue *Queue[T]) Dequeue() T {
	index := 0
	front := (*queue)[index]
	if len(*queue) > 1 {
		*queue = (*queue)[:index+1]
	} else {
		*queue = []T{}
	}
	return front
}
