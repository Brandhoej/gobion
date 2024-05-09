package algorithms

import "github.com/Brandhoej/gobion/pkg/structures"

func DFS[T any](
	successors func(node T) []T,
	visited func(state T) bool,
	yield func(node structures.LinkedNode[T]) bool,
	roots ...T,
) {
	stack := make(structures.Stack[structures.LinkedNode[T]], 0)

	for _, root := range roots {
		node := structures.NewLinkedNode(nil, root)
		stack.Push(node)
	}

	for !stack.IsEmpty() {
		parent := stack.Pop()

		for _, successor := range successors(parent.Data) {
			if !visited(successor) {
				node := structures.NewLinkedNode(&parent, successor)
				stack.Push(node)

				if !yield(node) {
					return
				}
			}
		}
	}
}
