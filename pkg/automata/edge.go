package automata

import (
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Edge struct {
	source      symbols.Symbol
	guard       Guard
	update      Update
	destination symbols.Symbol
}

func NewEdge(
	source symbols.Symbol, guard Guard, update Update, destination symbols.Symbol,
) Edge {
	return Edge{
		source:      source,
		guard:       guard,
		update:      update,
		destination: destination,
	}
}

func (edge Edge) Source() symbols.Symbol {
	return edge.source
}

func (edge Edge) Destination() symbols.Symbol {
	return edge.destination
}

func (edge Edge) IsEnabled(valuations language.Valuations, solver *Interpreter) bool {
	return edge.guard.IsSatisfied(valuations, solver)
}

func (edge Edge) Traverse(state State, solver *Interpreter) State {
	valuations := state.valuations.Copy()
	return NewState(
		edge.destination,
		valuations,
		edge.update.Apply(state.constraint),
	)
}
