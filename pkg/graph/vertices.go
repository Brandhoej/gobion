package graph

type Vertex any

type Vertices[K comparable, V Vertex] interface {
	Add(vertex V, key K) K
	Vertex(key K) (V, bool)
	All(yield func(K, V) bool)
}

type VertexMap[K comparable, V Vertex] struct {
	vertices map[K]V
}

func (mapping *VertexMap[K, V]) Add(vertex V, key K) K {
	mapping.vertices[key] = vertex
	return key
}

func (mapping *VertexMap[K, V]) Vertex(key K) (V, bool) {
	vertex, exists := mapping.vertices[key]
	return vertex, exists
}

func (mapping *VertexMap[K, V]) All(yield func(K, V) bool) {
	for key, vertex := range mapping.vertices {
		if !yield(key, vertex) {
			return
		}
	}
}

func NewVertexMap[K comparable, V Vertex]() *VertexMap[K, V] {
	return &VertexMap[K, V]{
		vertices: map[K]V{},
	}
}
