package structures

type Stack[T any] []T

func (stack *Stack[T]) IsEmpty() bool {
	return len(*stack) == 0
}

func (stack *Stack[T]) Push(elements ...T) {
	*stack = append(*stack, elements...)
}

func (stack *Stack[T]) Peek() (element T, exists bool) {
	if len(*stack) == 0 {
		return element, exists
	}

	index := len(*stack) - 1
	return (*stack)[index], true
}

func (stack *Stack[T]) Pop() T {
	index := len(*stack) - 1
	top := (*stack)[index]
	*stack = (*stack)[:index]
	return top
}
