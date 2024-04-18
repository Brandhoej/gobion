package automata

import (
	"fmt"
	"io"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/graph"
)

type Automaton struct {
	context *z3.Context
	graph   *graph.LabeledDirected[Edge, Location]
	initial graph.Key
}

func NewAutomaton(context *z3.Context, graph *graph.LabeledDirected[Edge, Location], initial graph.Key) *Automaton {
	return &Automaton{
		context: context,
		graph:   graph,
		initial: initial,
	}
}

func (automaton *Automaton) Initial() graph.Key {
	return automaton.initial
}

func (automaton *Automaton) Location(key graph.Key) (location Location, exists bool) {
	location, exists = automaton.graph.At(key)
	return location, exists
}

func (automaton *Automaton) Locations(yield func(graph.Key, Location) bool) {
	automaton.graph.Vertices(yield)
}

func (automaton *Automaton) Outgoing(location graph.Key) (edges []Edge) {
	return automaton.graph.From(location)
}

func (automaton *Automaton) Ingoing(location graph.Key) (edges []Edge) {
	return automaton.graph.To(location)
}

func (automaton *Automaton) Complete(variables language.Variables, complete func(graph.Key, Guard) graph.Key) {
	automaton.Locations(func(source graph.Key, location Location) bool {
		// The disjunction is the disjunction of all guards.
		var disjunction Guard
		if outgoings := automaton.Outgoing(source); len(outgoings) > 0 {
			guards := make([]Guard, len(outgoings))
			for idx := range outgoings {
				guards[idx] = outgoings[idx].guard
			}

			disjunction = guards[0]
			if len(guards) > 1 {
				disjunction = guards[0].Disjunction(guards[1:]...)
			}
		} else {
			// If there are not outgoing edges then we assume a false edge.
			disjunction = NewGuard(language.NewFalse())
		}

		// Constrain by the location's invariant.
		invariant := NewGuard(location.invariant.condition)
		missing := disjunction.Negation().Conjunction(invariant)

		// If the negation of all didisjunctionsjoined edge guards constrained by the invariant
		// still has a solution then we have a "missing" edge to the completion destination.
		if missing.IsSatisfiable(newSolver(nil, nil, nil)) {
			destination := complete(source, missing)
			update := NewUpdate()
			edge := NewEdge(source, missing, update, destination)
			automaton.graph.AddEdge(edge)
		}

		return true
	})
}

func (automaton *Automaton) DOT(writer io.Writer) {
	automaton.graph.DOT(
		writer,
		func(location Location) string {
			return fmt.Sprintf("%v", location.name)
		},
		func(edge Edge) string {
			guard := edge.guard.String()
			update := edge.update.String()
			return fmt.Sprintf("\n%s\n%s", guard, update)
		},
	)
}
