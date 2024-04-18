package graph

type Edge interface {
	Source() Key
	Destination() Key
}

type Edges[E Edge] interface {
	To(destination Key) (edges []E)
	From(source Key) (edges []E)
	Add(edge E)
	All(yield func(E) bool)
}

type EdgesMap[E Edge] struct {
	edges []E
	outgoings map[Key][]E
	ingoings map[Key][]E
}

func (mapping *EdgesMap[E]) Add(edge E) {
	mapping.edges = append(mapping.edges, edge)
	if outgoings, exists := mapping.outgoings[edge.Source()]; exists {
		mapping.outgoings[edge.Source()] = append(outgoings, edge)
	} else {
		mapping.outgoings[edge.Source()] = []E{ edge }
	}

	if ingoings, exists := mapping.ingoings[edge.Destination()]; exists {
		mapping.ingoings[edge.Destination()] = append(ingoings, edge)
	} else {
		mapping.ingoings[edge.Destination()] = []E{ edge }
	}
}

func (mapping *EdgesMap[E]) All(yield func(E) bool) {
	for _, edge := range mapping.edges {
		if !yield(edge) {
			return
		}
	}
}

func (mapping *EdgesMap[E]) To(destination Key) (edges []E) {
	return mapping.ingoings[destination]
}

func (mapping *EdgesMap[E]) From(source Key) (edges []E) {
	return mapping.outgoings[source]
}

func NewEdgesMap[E Edge]() *EdgesMap[E] {
	return &EdgesMap[E]{
		edges: make([]E, 0),
		outgoings: map[Key][]E{},
		ingoings: map[Key][]E{},
	}
}