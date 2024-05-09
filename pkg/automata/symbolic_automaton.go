package automata

import (
	"fmt"
	"io"

	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type SymbolicAutomaton struct {
	Automaton[Location, Edge]
}

func NewSymbolicAutomaton(automaton Automaton[Location, Edge]) *SymbolicAutomaton {
	return &SymbolicAutomaton{
		Automaton: automaton,
	}
}

func (automaton SymbolicAutomaton) Complete(interpreter *Interpreter, complete func(symbols.Symbol, Guard) symbols.Symbol) {
	automaton.Locations(func(source symbols.Symbol, location Location) bool {
		// The disjunction is the disjunction of all guard conditions.
		// If there are not outgoing edges then we assume a false edge condition.
		// This means that if there are no edge a true condition is the missing.
		var disjunction Guard
		if outgoings := automaton.Outgoing(source); len(outgoings) > 0 {
			guards := make([]Guard, len(outgoings))
			for idx := range outgoings {
				guards[idx] = outgoings[idx].guard
			}

			disjunction = guards[0]
			if len(guards) > 1 {
				// We do "guards[0]." and not "disjunction." to not have an additional "depth".
				disjunction = guards[0].Disjunction(guards[1:]...)
			}
		} else {
			disjunction = NewGuard(language.NewFalse())
		}

		// This is the missing guard condition.
		negation := disjunction.Negation()

		// Constrain by the location's invariant.
		invariant := NewGuard(location.invariant.condition)
		missing := negation.Conjunction(invariant)

		// If the negation of all disjunctions of edge guards constrained by the invariant
		// still has a solution then we have a "missing" edge to the completion destination.
		if missing.IsSatisfiable(interpreter) {
			destination := complete(source, negation)
			update := NewUpdate(language.NewTrue())
			edge := NewEdge(source, negation, update, destination)
			automaton.graph.AddEdge(edge)
		}

		return true
	})
}

func (automaton SymbolicAutomaton) DOT(writer io.Writer, symbols symbols.Store[any]) {
	automaton.graph.DOT(
		writer,
		func(location Location) string {
			return fmt.Sprintf("%s\n%s", location.name, location.invariant.String(symbols))
		},
		func(edge Edge) string {
			guard := edge.guard.String(symbols)
			update := edge.update.String(symbols)
			return fmt.Sprintf("%s\n%s", guard, update)
		},
	)
}
