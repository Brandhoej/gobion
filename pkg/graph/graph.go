package graph

import (
	"fmt"
	"io"
)

type Key int

type LabeledDirected[E Edge, V Vertex] struct {
	vertices Vertices[V]
	edges    Edges[E]
}

func NewLabeledDirected[E Edge, V Vertex](
	vertices Vertices[V], edges Edges[E],
) *LabeledDirected[E, V] {
	return &LabeledDirected[E, V]{
		vertices: vertices,
		edges:    edges,
	}
}

func (graph *LabeledDirected[E, V]) To(destination Key) (edges []E) {
	return graph.edges.To(destination)
}

func (graph *LabeledDirected[E, V]) From(source Key) (edges []E) {
	return graph.edges.From(source)
}

func (graph *LabeledDirected[E, V]) At(key Key) (V, bool) {
	vertex, exists := graph.vertices.Get(key)
	return vertex, exists
}

func (graph *LabeledDirected[E, V]) Vertices(yield func(Key, V) bool) {
	graph.vertices.All(yield)
}

func (graph *LabeledDirected[E, V]) AddVertex(vertex V) Key {
	return graph.vertices.Add(vertex)
}

func (graph *LabeledDirected[E, V]) AddEdge(edge E) {
	graph.edges.Add(edge)
}

func (graph *LabeledDirected[E, V]) DOT(
	writer io.Writer,
	vertexLabel func(V) string,
	edgeLabel func(E) string,
) {
	io.WriteString(writer, "digraph G {\n")
	graph.vertices.All(func(key Key, vertex V) bool {
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