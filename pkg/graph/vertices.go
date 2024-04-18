package graph

type Vertex any

type Vertices[V Vertex] interface {
	Add(vertex V) Key
	Get(key Key) (V, bool)
	All(yield func(Key, V) bool)
}

type VertexMap[V Vertex] struct {
	vertices map[Key]V
}

func (mapping *VertexMap[V]) Add(vertex V) Key {
	key := Key(len(mapping.vertices))
	mapping.vertices[key] = vertex
	return key
}

func (mapping *VertexMap[V]) Get(key Key) (V, bool) {
	vertex, exists := mapping.vertices[key]
	return vertex, exists
}

func (mapping *VertexMap[V]) All(yield func(Key, V) bool) {
	for key, vertex := range mapping.vertices {
		if !yield(key, vertex) {
			return
		}
	}
}

func NewVertexMap[V Vertex]() *VertexMap[V] {
	return &VertexMap[V]{
		vertices: map[Key]V{},
	}
}