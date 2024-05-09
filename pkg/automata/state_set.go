package automata

import "github.com/Brandhoej/gobion/pkg/symbols"

type StateSet struct {
	states map[symbols.Symbol][]State
}

func NewStateSet() StateSet {
	return StateSet{
		states: make(map[symbols.Symbol][]State),
	}
}

func (set StateSet) Insert(solver *Interpreter, states ...State) (counter int) {
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

func (set StateSet) Contains(target State, solver *Interpreter) bool {
	if states, exists := set.states[target.location]; exists {
		for _, state := range states {
			if target.SubsetOf(state, solver) {
				return true
			}
		}
	}

	return false
}
