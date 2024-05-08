package language

import "github.com/Brandhoej/gobion/pkg/symbols"

type Valuations interface {
	Assign(symbol symbols.Symbol, value Expression)
	Value(symbol symbols.Symbol) (value Expression, exists bool)
	All(yield func(symbol symbols.Symbol, value Expression) bool) bool
	Copy() Valuations
}

type ValuationsMap struct {
	valuations map[symbols.Symbol]Expression
}

func NewValuationsMap() *ValuationsMap {
	return &ValuationsMap{
		valuations: map[symbols.Symbol]Expression{},
	}
}

func (mapping *ValuationsMap) Value(symbol symbols.Symbol) (value Expression, exists bool) {
	value, exists = mapping.valuations[symbol]
	return value, exists
}

func (mapping *ValuationsMap) Assign(symbol symbols.Symbol, expression Expression) {
	mapping.valuations[symbol] = expression
}

func (mapping *ValuationsMap) All(
	yield func(symbol symbols.Symbol, value Expression) bool,
) bool {
	for symbol, value := range mapping.valuations {
		if !yield(symbol, value) {
			return false
		}
	}
	return true
}

func (mapping *ValuationsMap) Copy() Valuations {
	copy := NewValuationsMap()
	mapping.All(func(symbol symbols.Symbol, value Expression) bool {
		copy.Assign(symbol, value)
		return true
	})
	return copy
}