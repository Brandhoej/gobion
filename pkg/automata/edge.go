package automata

import "github.com/Brandhoej/gobion/pkg/graph"

type EdgeDirection bool

type Edge struct {
	source      graph.Key
	guard       Guard
	update      Update
	destination graph.Key
}

func NewEdge(
	source graph.Key, guard Guard, update Update, destination graph.Key,
) Edge {
	return Edge{
		source:      source,
		guard:       guard,
		update:      update,
		destination: destination,
	}
}

func (edge Edge) Source() graph.Key {
	return edge.source
}

func (edge Edge) Destination() graph.Key {
	return edge.destination
}

func (edge Edge) IsEnabled(variables Variables, state State) bool {
	return edge.guard.IsSatisfiable(variables, state.valuations)
}

func (edge Edge) Traverse(variables Variables, state State, forward bool) State {
	valuations := edge.update.Apply(variables, state.valuations, forward)
	return NewState(edge.destination, valuations)
}
