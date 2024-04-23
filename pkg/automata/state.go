package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/graph"
)

type State struct {
	location   graph.Key
	valuations expressions.Valuations[*z3.AST]
	constraint constraints.Constraint
}

func NewState(
	location graph.Key,
	valuations expressions.Valuations[*z3.AST],
	constraint constraints.Constraint,
) State {
	return State{
		location:   location,
		valuations: valuations,
		constraint: constraint,
	}
}

func (state State) SubsetOf(other State, solver *ConstraintSolver) bool {
	if state.location != other.location {
		return false
	}
	return solver.Satisfies(
		state.valuations,
		constraints.Implication(
			constraints.Conjunction(state.constraint, other.constraint),
			state.constraint,
		),
	)
}
