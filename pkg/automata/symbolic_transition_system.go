package automata

import "github.com/Brandhoej/gobion/pkg/automata/language"

type SymbolicTransitionSystem struct {
	automaton   *SymbolicAutomaton
	interpreter *Interpreter
}

func NewTransitionSystem(
	automaton *SymbolicAutomaton, interpreter *Interpreter,
) *SymbolicTransitionSystem {
	return &SymbolicTransitionSystem{
		automaton: automaton,
		interpreter: interpreter,
	}
}

func (system *SymbolicTransitionSystem) Initial(valuations language.Valuations) State {
	initial := system.automaton.initial
	location, _ := system.automaton.Location(initial)
	return NewState(initial, valuations, location.invariant.condition)
}

// Returns all states from the state.
func (system *SymbolicTransitionSystem) Outgoing(state State) (successors []State) {
	if location, exists := system.automaton.Location(state.location); exists {
		// We have found an inconsistency where the location is disabled.
		// Meaning that even enabled edges wont be traversable.
		if !location.IsEnabled(state.valuations, system.interpreter) {
			return successors
		}
	} else {
		panic("State is in an unkown location")
	}

	edges := system.automaton.Outgoing(state.location)
	for _, edge := range edges {
		// Check if we can even traverse the edge.
		if !edge.IsEnabled(state.valuations, system.interpreter) {
			continue
		}

		// We can traverse the edge so we create a new and updated state.
		state := edge.Traverse(state, system.interpreter)
		successors = append(successors, state)
	}
	return successors
}

func (system *SymbolicTransitionSystem) Reachability(
	valuations language.Valuations, search SearchStrategy, goals ...State,
) Trace {
	return search.For(
		func(state State) bool {
			// We have reached a goal when the locations are the same
			// and the goal contains (Meaning that more valuations or the same) are possible.
			for _, goal := range goals {
				if goal.SubsetOf(state, system.interpreter) {
					return true
				}
			}

			return false
		},
		system.Initial(valuations),
	)
}
