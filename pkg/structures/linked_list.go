package structures

type LinkedNode[T any] struct {
	Parent *LinkedNode[T]
	Data T
}

func NewLinkedNode[T any](parent *LinkedNode[T], data T) LinkedNode[T] {
	return LinkedNode[T]{
		Parent: parent,
		Data: data,
	}
}

func (node LinkedNode[T]) Array() (data []T) {
	for current := &node; current != nil; current = current.Parent {
		data = append(data, current.Data)
	}
	return data
}