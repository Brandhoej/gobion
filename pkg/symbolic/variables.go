package symbolic

import (
	"fmt"

	"github.com/Brandhoej/gobion/internal/z3"
)

type Variables interface {
	Declare(symbol Symbol, identifier string, sort *z3.Sort) (variable *z3.AST)
	Advance(symbol Symbol) (variable *z3.AST, exists bool)
	Variable(symbol Symbol) (variable *z3.AST, exists bool)
}

type VariablesMap struct {
	identifiers map[Symbol]string
	counters    map[Symbol]int
	variables   map[Symbol]*z3.AST
}

func NewVariablesMap() *VariablesMap {
	return &VariablesMap{
		identifiers: map[Symbol]string{},
		counters:    map[Symbol]int{},
		variables:   map[Symbol]*z3.AST{},
	}
}

func (mapping *VariablesMap) Declare(symbol Symbol, identifier string, sort *z3.Sort) (variable *z3.AST) {
	mapping.identifiers[symbol] = identifier
	mapping.counters[symbol] = 0
	variable = sort.Context().NewConstant(
		z3.WithName(identifier+"_0"), sort,
	)
	mapping.variables[symbol] = variable
	return
}

func (mapping *VariablesMap) Advance(symbol Symbol) (variable *z3.AST, exists bool) {
	variable, exists = mapping.Variable(symbol)
	if !exists {
		return
	}

	original := mapping.identifiers[symbol]
	counter := mapping.counters[symbol] + 1
	variable = variable.Context().NewConstant(
		z3.WithName(fmt.Sprintf("%s_%v", original, counter)), variable.Sort(),
	)
	mapping.variables[symbol] = variable
	mapping.counters[symbol] = counter

	return
}

func (mapping *VariablesMap) Variable(symbol Symbol) (variable *z3.AST, exists bool) {
	variable, exists = mapping.variables[symbol]
	return
}
