package automata

import (
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type State struct {
	location   symbols.Symbol
	valuations language.Valuations
	constraint language.Expression
}

func NewState(
	location symbols.Symbol,
	valuations language.Valuations,
	constraint language.Expression,
) State {
	return State{
		location:   location,
		valuations: valuations,
		constraint: constraint,
	}
}

func (state State) SubsetOf(other State, interpreter *Interpreter) bool {
	if state.location != other.location {
		return false
	}
	return interpreter.IsSatisfied(
		state.valuations,
		language.NewBinary(
			language.Conjunction(state.constraint, other.constraint),
			language.Implication,
			state.constraint,
		),
	)
}
