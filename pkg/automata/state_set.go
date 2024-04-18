package automata

import "github.com/Brandhoej/gobion/pkg/graph"

type StateSet struct {
	states map[graph.Key][]State
}

func NewStateSet(states ...State) StateSet {
	set := StateSet{
		states: make(map[graph.Key][]State),
	}
	for _, state := range states {
		set.Insert(state)
	}
	return set
}

func (set StateSet) Insert(state State) {
	if states, exists := set.states[state.location]; exists {
		if set.Contains(state) {
			return
		}

		set.states[state.location] = append(states, state)
	} else {
		set.states[state.location] = []State{state}
	}
}

func (set StateSet) Contains(target State) bool {
	if states, exists := set.states[target.location]; exists {
		for _, state := range states {
			if target.valuations.SubsetOf(state.valuations) {
				return true
			}
		}
	}

	return false
}