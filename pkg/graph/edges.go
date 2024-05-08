package graph

type Edge[K comparable] interface {
	Source() K
	Destination() K
}

type Edges[K comparable, E Edge[K]] interface {
	To(destination K) (edges []E)
	From(source K) (edges []E)
	Connect(edge E)
	All(yield func(E) bool)
}

type EdgeSlice[K comparable, E Edge[K]] []E

func (slice *EdgeSlice[K, E]) To(destination K) (edges []E) {
	for idx := range *slice {
		if (*slice)[idx].Destination() == destination {
			edges = append(edges, (*slice)[idx])
		}
	}

	return edges
}

func (slice *EdgeSlice[K, E]) From(source K) (edges []E) {
	for idx := range *slice {
		if (*slice)[idx].Source() == source {
			edges = append(edges, (*slice)[idx])
		}
	}

	return edges
}

func (slice *EdgeSlice[K, E]) Connect(edge E) {
	(*slice) = append((*slice), edge)
}

func (slice *EdgeSlice[K, E]) All(yield func(E) bool) {
	for idx := range *slice {
		if !yield((*slice)[idx]) {
			return
		}
	}
}

type EdgeMap[K comparable, E Edge[K]] struct {
	edges     []E
	outgoings map[K][]E
	ingoings  map[K][]E
}

func NewEdgesMap[K comparable, E Edge[K]]() *EdgeMap[K, E] {
	return &EdgeMap[K, E]{
		edges:     make([]E, 0),
		outgoings: map[K][]E{},
		ingoings:  map[K][]E{},
	}
}

func (mapping *EdgeMap[K, E]) Connect(edge E) {
	mapping.edges = append(mapping.edges, edge)
	if outgoings, exists := mapping.outgoings[edge.Source()]; exists {
		mapping.outgoings[edge.Source()] = append(outgoings, edge)
	} else {
		mapping.outgoings[edge.Source()] = []E{edge}
	}

	if ingoings, exists := mapping.ingoings[edge.Destination()]; exists {
		mapping.ingoings[edge.Destination()] = append(ingoings, edge)
	} else {
		mapping.ingoings[edge.Destination()] = []E{edge}
	}
}

func (mapping *EdgeMap[K, E]) All(yield func(E) bool) {
	for _, edge := range mapping.edges {
		if !yield(edge) {
			return
		}
	}
}

func (mapping *EdgeMap[K, E]) To(destination K) (edges []E) {
	return mapping.ingoings[destination]
}

func (mapping *EdgeMap[K, E]) From(source K) (edges []E) {
	return mapping.outgoings[source]
}