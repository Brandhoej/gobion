package automata

import (
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/graph"
)

type State struct {
	location   graph.Key
	valuations language.Valuations
	constraint language.Expression
}

func NewState(
	location graph.Key,
	valuations language.Valuations,
	constraint language.Expression,
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
		language.NewBinary(
			language.Conjunction(state.constraint, other.constraint),
			language.Implication,
			state.constraint,
		),
	)
}
