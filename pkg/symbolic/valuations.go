package symbolic

import "github.com/Brandhoej/gobion/internal/z3"

type Valuations interface {
	Load(symbol Symbol) (valuation *z3.AST)
	Store(symbol Symbol, value *z3.AST)
	Symbols(function func(symbol Symbol))
}

type ValuationsMap struct {
	valuations map[Symbol]*z3.AST
}

func NewValuationsMap() *ValuationsMap {
	return &ValuationsMap{
		valuations: map[Symbol]*z3.AST{},
	}
}

func (mapping *ValuationsMap) Store(symbol Symbol, value *z3.AST) {
	mapping.valuations[symbol] = value
}

func (mapping *ValuationsMap) Load(symbol Symbol) (valuation *z3.AST) {
	valuation = mapping.valuations[symbol]
	return
}

func (mapping *ValuationsMap) Symbols(function func(symbol Symbol)) {
	for symbol := range mapping.valuations {
		function(symbol)
	}
}
