package automata

import (
	"fmt"
	"io"
	"slices"

	"github.com/Brandhoej/gobion/pkg/symbols"
)

type IOAutomaton struct {
	Automaton[Location, IOEdge]
	inputs  []Action
	outputs []Action
}

func NewIOAutomaton(automaton Automaton[Location, IOEdge], inputs, outputs []Action) *IOAutomaton {
	return &IOAutomaton{
		Automaton: automaton,
		inputs:    inputs,
		outputs:   outputs,
	}
}

func (automaton IOAutomaton) IsInput(action Action) bool {
	return slices.Contains(automaton.inputs, action)
}

func (automaton IOAutomaton) IsOutput(action Action) bool {
	return slices.Contains(automaton.outputs, action)
}

// Returns all edges to the location (Ingoing) where its action is the action.
func (automaton IOAutomaton) Ingoing(location symbols.Symbol, action Action) []IOEdge {
	ingoing := automaton.Automaton.Ingoing(location)
	edges := make([]IOEdge, 0, len(ingoing))
	for idx := range ingoing {
		if ingoing[idx].action == action {
			edges = append(edges, ingoing[idx])
		}
	}
	return edges
}

// Returns all edges from the location (Outgoing) where its action is the action.
func (automaton IOAutomaton) Outgoing(location symbols.Symbol, action Action) []IOEdge {
	outgoing := automaton.Automaton.Outgoing(location)
	edges := make([]IOEdge, 0, len(outgoing))
	for idx := range outgoing {
		if outgoing[idx].action == action {
			edges = append(edges, outgoing[idx])
		}
	}
	return edges
}

func (automaton *IOAutomaton) DOT(writer io.Writer, store symbols.Store[any]) {
	automaton.graph.DOT(
		writer,
		func(location Location) string {
			return fmt.Sprintf("%s\n%s", location.name, location.invariant.String(store))
		},
		func(edge IOEdge) string {
			guard := edge.guard.String(store)
			update := edge.update.String(store)
			action := edge.action.String(
				automaton.IsInput(edge.action), store,
			)
			return fmt.Sprintf("%s\n%s\n%s", guard, action, update)
		},
	)
}
