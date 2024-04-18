package automata

import (
	"github.com/Brandhoej/gobion/pkg/graph"
)

type State struct {
	location graph.Key
	valuations Valuations
}

func NewState(
	location graph.Key,
	valuations Valuations,
) State {
	return State{
		location: location,
		valuations: valuations,
	}
}