package automata

import (
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/graph"
)

type State struct {
	location graph.Key
	valuations language.Valuations
}

func NewState(
	location graph.Key,
	valuations language.Valuations,
) State {
	return State{
		location: location,
		valuations: valuations,
	}
}

func (state State) SubsetOf(other State) bool {
	if state.location == other.location {
		return true
	}
	// TODO: Check valuations.
	panic("Not implemented yet")
}