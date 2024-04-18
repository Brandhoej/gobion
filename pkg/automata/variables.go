package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Variables interface {
	Declare(symbol symbols.Symbol, sort *z3.Sort) (variable *z3.AST)
	Variable(symbol symbols.Symbol) (variable *z3.AST, exists bool)
	All(yield func(symbol symbols.Symbol, variable *z3.AST) bool) bool
}

type VariablesMap struct {
	context *z3.Context
	variables map[symbols.Symbol]*z3.AST
}

func NewVariablesMap(context *z3.Context) *VariablesMap {
	return &VariablesMap{
		context: context,
		variables: map[symbols.Symbol]*z3.AST{},
	}
}

func (mapping *VariablesMap) Declare(symbol symbols.Symbol, sort *z3.Sort) (variable *z3.AST) {
	if variable, exists := mapping.Variable(symbol); exists {
		return variable
	}
	variable = mapping.context.NewConstant(
		z3.WithInt(int(symbol)), sort,
	)
	mapping.variables[symbol] = variable
	return variable
}

func (mapping *VariablesMap) Variable(symbol symbols.Symbol) (variable *z3.AST, exists bool) {
	variable, exists = mapping.variables[symbol]
	return variable, exists
}

func (mapping *VariablesMap) All(
	yield func(symbol symbols.Symbol, variable *z3.AST) bool,
) bool {
	for symbol, variable := range mapping.variables {
		if !yield(symbol, variable) {
			return false
		}
	}
	return true
}