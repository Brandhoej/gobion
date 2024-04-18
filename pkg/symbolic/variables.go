package symbolic

import "github.com/Brandhoej/gobion/internal/z3"

type Variables interface {
	Declare(symbol Symbol, sort *z3.Sort) (variable *z3.AST)
	Get(symbol Symbol) (variable *z3.AST)
	Symbols(function func(symbol Symbol))
}

type VariablesMap struct {
	variables map[Symbol]*z3.AST
}

func NewVariablesMap() *VariablesMap {
	return &VariablesMap{
		variables: map[Symbol]*z3.AST{},
	}
}

func (mapping *VariablesMap) Declare(symbol Symbol, sort *z3.Sort) (variable *z3.AST) {
	context := sort.Context()
	variable = context.NewConstant(
		z3.WithSymbol(symbol.Z3(context)), sort,
	)
	mapping.variables[symbol] = variable
	return variable
}

func (mapping *VariablesMap) Get(symbol Symbol) (variable *z3.AST) {
	variable = mapping.variables[symbol]
	return variable
}

func (mapping *VariablesMap) Symbols(function func(symbol Symbol)) {
	for symbol := range mapping.variables {
		function(symbol)
	}
}
