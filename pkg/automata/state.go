package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/graph"
)

type State struct {
	location   graph.Key
	valuations expressions.Valuations[*z3.AST]
	constraint expressions.Expression
}

func NewState(
	location graph.Key,
	valuations expressions.Valuations[*z3.AST],
	constraint expressions.Expression,
) State {
	return State{
		location:   location,
		valuations: valuations,
		constraint: constraint,
	}
}

func (state State) SubsetOf(other State, solver *Interpreter) bool {
	if state.location != other.location {
		return false
	}
	return solver.IsSatisfied(
		state.valuations,
		expressions.NewBinary(
			expressions.Conjunction(state.constraint, other.constraint),
			expressions.Implication,
			state.constraint,
		),
	)
}
