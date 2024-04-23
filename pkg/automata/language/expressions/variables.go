package expressions

import "github.com/Brandhoej/gobion/pkg/symbols"

type Variables[T any] interface {
	Declare(symbol symbols.Symbol, variable T)
	Lookup(symbol symbols.Symbol) (variable T, exists bool)
	All(yield func(symbol symbols.Symbol, variable T) bool) bool
}

type VariablesMap[T any] struct {
	variables map[symbols.Symbol]T
}

func NewVariablesMap[T any]() *VariablesMap[T] {
	return &VariablesMap[T]{
		variables: map[symbols.Symbol]T{},
	}
}

func (mapping *VariablesMap[T]) Declare(symbol symbols.Symbol, variable T) {
	mapping.variables[symbol] = variable
}

func (mapping *VariablesMap[T]) Lookup(symbol symbols.Symbol) (variable T, exists bool) {
	variable, exists = mapping.variables[symbol]
	return variable, exists
}

func (mapping *VariablesMap[T]) All(yield func(symbol symbols.Symbol, variable T) bool) bool {
	for symbol, sort := range mapping.variables {
		if !yield(symbol, sort) {
			return false
		}
	}
	return true
}
