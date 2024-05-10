package automata

import (
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/Brandhoej/gobion/pkg/zones"
)

type TIOAutomaton struct {
	IOAutomaton
	clocks symbols.Store[zones.Clock]
}

func NewTIOAutomaton(automaton IOAutomaton, clocks symbols.Store[zones.Clock]) *TIOAutomaton {
	return &TIOAutomaton{
		IOAutomaton: automaton,
		clocks: clocks,
	}
}