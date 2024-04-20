package automata

import "github.com/Brandhoej/gobion/pkg/graph"

type StateSet struct {
	states map[graph.Key][]State
}

func NewStateSet() StateSet {
	return StateSet{
		states: make(map[graph.Key][]State),
	}
}

func (set StateSet) Insert(solver *ConstraintSolver, states ...State) (counter int) {
	for _, state := range states {
		if states, exists := set.states[state.location]; exists {
			if set.Contains(state, solver) {
				continue
			}

			set.states[state.location] = append(states, state)
		} else {
			set.states[state.location] = []State{state}
		}

		counter += 1
	}

	return counter
}

func (set StateSet) Contains(target State, solver *ConstraintSolver) bool {
	if states, exists := set.states[target.location]; exists {
		for _, state := range states {
			if target.SubsetOf(state, solver) {
				return true
			}
		}
	}

	return false
}
