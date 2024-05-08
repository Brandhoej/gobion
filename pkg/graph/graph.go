package graph

import (
	"fmt"
	"io"
)

type LabeledDirected[K comparable, E Edge[K], V Vertex] struct {
	vertices Vertices[K, V]
	edges    Edges[K, E]
}

func NewLabeledDirected[K comparable, E Edge[K], V Vertex](
	vertices Vertices[K, V], edges Edges[K, E],
) *LabeledDirected[K, E, V] {
	return &LabeledDirected[K, E, V]{
		vertices: vertices,
		edges:    edges,
	}
}

func (graph *LabeledDirected[K, E, V]) To(destination K) (edges []E) {
	return graph.edges.To(destination)
}

func (graph *LabeledDirected[K, E, V]) From(source K) (edges []E) {
	return graph.edges.From(source)
}

func (graph *LabeledDirected[K, E, V]) At(key K) (V, bool) {
	vertex, exists := graph.vertices.Vertex(key)
	return vertex, exists
}

func (graph *LabeledDirected[K, E, V]) Vertices(yield func(K, V) bool) {
	graph.vertices.All(yield)
}

func (graph *LabeledDirected[K, E, V]) AddVertex(vertex V, key K) K {
	return graph.vertices.Add(vertex, key)
}

func (graph *LabeledDirected[K, E, V]) AddEdge(edge E) {
	graph.edges.Connect(edge)
}

func (graph *LabeledDirected[K, E, V]) DOT(
	writer io.Writer,
	vertexLabel func(V) string,
	edgeLabel func(E) string,
) {
	io.WriteString(writer, "digraph G {\n")
	graph.vertices.All(func(key K, vertex V) bool {
		str := fmt.Sprintf(
			"%v [label=\"%s\"]\n",
			key,
			vertexLabel(vertex),
		)
		io.WriteString(writer, str)
		return true
	})
	graph.edges.All(func(edge E) bool {
		str := fmt.Sprintf(
			"%v -> %v [label=\"%s\"]\n",
			edge.Source(),
			edge.Destination(),
			edgeLabel(edge),
		)
		io.WriteString(writer, str)
		return true
	})
	io.WriteString(writer, "}\n")
}
