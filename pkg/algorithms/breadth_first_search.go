package algorithms

import "github.com/Brandhoej/gobion/pkg/structures"

func BFS[T any](
	successors func(node T) []T,
	visited func(state T) bool,
	yield func(node structures.LinkedNode[T]) bool,
	roots ...T,
) {
	queue := make(structures.Queue[structures.LinkedNode[T]], 0)
	
	for _, root := range roots {
		node := structures.NewLinkedNode(nil, root)
		queue.Enqueue(node)
	}

	for !queue.IsEmpty() {
		parent := queue.Dequeue()

		for _, successor := range successors(parent.Data) {			
			if !visited(successor) {
				node := structures.NewLinkedNode(&parent, successor)
				queue.Enqueue(node)
				
				if !yield(node) {
					return
				}
			}
		}
	}
}