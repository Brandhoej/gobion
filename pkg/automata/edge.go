package automata

import (
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
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

func (edge Edge) IsEnabled(valuations expressions.Valuations, solver *ConstraintSolver) bool {
	return edge.guard.IsSatisfiable(valuations, solver)
}

func (edge Edge) Traverse(state State, solver *ConstraintSolver) State {
	valuations := state.valuations.Copy()
	edge.update.Apply(valuations, solver)
	return NewState(
		edge.destination,
		valuations,
		constraints.Conjunction(state.constraint, edge.update.constraint),
	)
}
