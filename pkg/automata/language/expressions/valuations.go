package expressions

import (
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Valuations[T any] interface {
	Assign(symbol symbols.Symbol, value T)
	Value(symbol symbols.Symbol) (value T, exists bool)
	All(yield func(symbol symbols.Symbol, value T) bool) bool
	Copy() Valuations[T]
}

type ValuationsMap[T any] struct {
	valuations map[symbols.Symbol]T
}

func NewValuationsMap[T any]() *ValuationsMap[T] {
	return &ValuationsMap[T]{
		valuations: map[symbols.Symbol]T{},
	}
}

func (mapping *ValuationsMap[T]) Value(symbol symbols.Symbol) (value T, exists bool) {
	value, exists = mapping.valuations[symbol]
	return value, exists
}

func (mapping *ValuationsMap[T]) Assign(symbol symbols.Symbol, expression T) {
	mapping.valuations[symbol] = expression
}

func (mapping *ValuationsMap[T]) All(
	yield func(symbol symbols.Symbol, value T) bool,
) bool {
	for symbol, value := range mapping.valuations {
		if !yield(symbol, value) {
			return false
		}
	}
	return true
}

func (mapping *ValuationsMap[T]) Copy() Valuations[T] {
	copy := NewValuationsMap[T]()
	mapping.All(func(symbol symbols.Symbol, value T) bool {
		copy.Assign(symbol, value)
		return true
	})
	return copy
}