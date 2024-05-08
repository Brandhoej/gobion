package automata

import (
	"github.com/Brandhoej/gobion/pkg/graph"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Automaton[L graph.Vertex, E graph.Edge[symbols.Symbol]] struct {
	graph   *graph.LabeledDirected[symbols.Symbol, E, L]
	initial symbols.Symbol
}

func NewAutomaton(
	graph *graph.LabeledDirected[symbols.Symbol, Edge, Location],
	initial symbols.Symbol,
) *Automaton[Location, Edge] {
	return &Automaton[Location, Edge]{
		graph:   graph,
		initial: initial,
	}
}

func (automaton *Automaton[L, E]) Initial() symbols.Symbol {
	return automaton.initial
}

func (automaton *Automaton[L, E]) Location(key symbols.Symbol) (location L, exists bool) {
	location, exists = automaton.graph.At(key)
	return location, exists
}

func (automaton *Automaton[L, E]) Locations(yield func(symbols.Symbol, L) bool) {
	automaton.graph.Vertices(yield)
}

func (automaton *Automaton[L, E]) Outgoing(location symbols.Symbol) (edges []E) {
	return automaton.graph.From(location)
}

func (automaton *Automaton[L, E]) Ingoing(location symbols.Symbol) (edges []E) {
	return automaton.graph.To(location)
}