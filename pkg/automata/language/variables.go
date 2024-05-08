package language

import "github.com/Brandhoej/gobion/pkg/symbols"

type Variables interface {
	Declare(symbol symbols.Symbol, sort Sort)
	Lookup(symbol symbols.Symbol) (sort Sort, exists bool)
	All(yield func(symbol symbols.Symbol, sort Sort) bool) bool
}

type VariablesMap struct {
	variables map[symbols.Symbol]Sort
}

func NewVariablesMap() *VariablesMap {
	return &VariablesMap{
		variables: map[symbols.Symbol]Sort{},
	}
}

func (mapping *VariablesMap) Declare(symbol symbols.Symbol, sort Sort) {
	mapping.variables[symbol] = sort
}

func (mapping *VariablesMap) Lookup(symbol symbols.Symbol) (sort Sort, exists bool) {
	sort, exists = mapping.variables[symbol]
	return sort, exists
}

func (mapping *VariablesMap) All(yield func(symbol symbols.Symbol, sort Sort) bool) bool {
	for symbol, sort := range mapping.variables {
		if !yield(symbol, sort) {
			return false
		}
	}
	return true
}
