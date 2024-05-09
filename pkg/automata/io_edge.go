package automata

import "github.com/Brandhoej/gobion/pkg/symbols"

type IOEdge struct {
	Edge
	action Action
}

func NewIOEdge(
	source symbols.Symbol, action Action, guard Guard, update Update, destination symbols.Symbol,
) IOEdge {
	return IOEdge{
		Edge:   NewEdge(source, guard, update, destination),
		action: action,
	}
}
