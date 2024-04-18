package automata

import (
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/graph"
)

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

func (edge Edge) IsEnabled(solver Solver) bool {
	return edge.guard.IsSatisfiable(solver)
}

func (edge Edge) Traverse(variables language.Variables, state State) State {
	/*valuations := edge.update.Apply(variables, state.valuations)
	return NewState(edge.destination, valuations)*/
	panic("Not implemented yet")
}
