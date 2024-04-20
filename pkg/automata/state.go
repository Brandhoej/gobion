package automata

import (
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/graph"
)

type State struct {
	location graph.Key
	constraint constraints.Constraint
}

func NewState(
	location graph.Key,
	constraint constraints.Constraint,
) State {
	return State{
		location: location,
		constraint: constraint,
	}
}

func (state State) SubsetOf(other State, solver ConstraintSolver) bool {
	if state.location == other.location {
		return true
	}
	return solver.HasSolutionFor(
		constraints.Conjunction(state.constraint, other.constraint),
	)
}