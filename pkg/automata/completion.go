package automata

import "github.com/Brandhoej/gobion/pkg/symbols"

var AngelicCompletion = func(location symbols.Symbol, _ Guard) symbols.Symbol {
	return location
}

var DirectedCompletion = func(destination symbols.Symbol) func(symbols.Symbol, Guard) symbols.Symbol {
	return func(symbols.Symbol, Guard) symbols.Symbol {
		return destination
	}
}