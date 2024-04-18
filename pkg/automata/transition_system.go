package automata

import "fmt"

type TransitionSystem struct {
	variables Variables
	automaton *Automaton
}

func NewTransitionSystem(variable Variables, automaton *Automaton) *TransitionSystem {
	return &TransitionSystem{
		variables: variable,
		automaton: automaton,
	}
}

func (system *TransitionSystem) Initial(valuations Valuations) State {
	return NewState(system.automaton.initial, valuations)
}

// Returns all states from the state.
func (system *TransitionSystem) Outgoing(state State) (successors []State) {
	if location, exists := system.automaton.Location(state.location); exists {
		// We have found an inconsistency where the location is disabled.
		// Meaning that even enabled edges wont be traversable.
		if !location.IsEnabled(system.variables, state.valuations) {
			return successors
		}
	} else {
		panic("State is in an unkown location")
	}

	edges := system.automaton.Outgoing(state.location)
	for _, edge := range edges {
		// Check if we can even traverse the edge.
		if !edge.IsEnabled(system.variables, state) {
			continue
		}

		// We can traverse the edge so we create a new and updated state.
		state := edge.Traverse(system.variables, state, true)
		successors = append(successors, state)
	}
	return successors
}

// Returns all states to the state.
func (system *TransitionSystem) Ingoing(state State) (predecessors []State) {
	if location, exists := system.automaton.Location(state.location); exists {
		// We have found an inconsistency where the location is disabled.
		// Meaning that even enabled edges wont be traversable.
		if !location.IsEnabled(system.variables, state.valuations) {
			return predecessors
		}
	} else {
		panic("State is in an unkown location")
	}

	edges := system.automaton.Ingoing(state.location)
	for _, edge := range edges {
		// Apply edge guards to variables as they must be satisfied
		// for the guard to be enabled in the first place.
		fmt.Println(edge)

		// Next, apply the invariant to the state, as this
		// must also be satisfied in the first place.

		predecessors = append(predecessors, state)
	}
	return predecessors
}

func (system *TransitionSystem) Reachability(search SearchStrategy, goals ...State) Trace {
	return search.For(
		func(state State) bool {
			// We have reached a goal when the locations are the same
			// and the goal contains (Meaning that more valuations or the same) are possible.
			for _, goal := range goals {
				if state.location == goal.location &&
					goal.valuations.SubsetOf(state.valuations) {
					return true
				}
			}

			return false
		},
		system.Initial(nil), // TODO
	)
}
