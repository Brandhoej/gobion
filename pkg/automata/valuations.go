package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Value struct {
	ast      *z3.AST
	relation Equality
}

func NewValue(ast *z3.AST, relation Equality) Value {
	return Value{
		ast:      ast,
		relation: relation,
	}
}

type Valuations interface {
	Value(symbol symbols.Symbol) (value Value, exists bool)
	Assign(symbol symbols.Symbol, value *z3.AST, operator Equality)
	All(yield func(symbol symbols.Symbol, value Value) bool) bool
	SubsetOf(others ...Valuations) bool
}

type ValuationsMap struct {
	context    *z3.Context
	valuations map[symbols.Symbol]Value
}

func NewValuationsMap(context *z3.Context) *ValuationsMap {
	return &ValuationsMap{
		context:    context,
		valuations: map[symbols.Symbol]Value{},
	}
}

func (mapping *ValuationsMap) Value(symbol symbols.Symbol) (value Value, exists bool) {
	value, exists = mapping.valuations[symbol]
	return value, exists
}

func (mapping *ValuationsMap) Assign(symbol symbols.Symbol, value *z3.AST, relation Equality) {
	mapping.valuations[symbol] = NewValue(value, relation)
}

func (mapping *ValuationsMap) Negation() (valuations Valuations) {
	valuations = NewValuationsMap(mapping.context)
	for symbol, value := range mapping.valuations {
		valuations.Assign(symbol, value.ast, value.relation.Negation())
	}
	return valuations
}

func (mapping *ValuationsMap) SubsetOf(others ...Valuations) bool {
	context := mapping.context
	solver := context.NewSolver()

	variables := make(map[symbols.Symbol]*z3.AST)
	valuations := make(map[symbols.Symbol]*z3.AST)

	conjoin := func(symbol symbols.Symbol, current Value) {
		variable, exists := variables[symbol]
		if !exists {
			variable = context.NewConstant(
				z3.WithInt(int(symbol)), current.ast.Sort(),
			)
			variables[symbol] = variable
		}

		assignment := current.relation.ast(variable, current.ast)
		if valuation, exists := valuations[symbol]; exists {
			valuations[symbol] = z3.And(valuation, assignment)
		} else {
			valuations[symbol] = assignment
		}
	}
	valuationsCounter := 0
	for _, other := range others {
		other.All(func(symbol symbols.Symbol, value Value) bool {
			conjoin(symbol, value)
			valuationsCounter += 1
			return true
		})

		// Pigeonhole principle ("other" must contain all symbols of "mapping").
		if valuationsCounter < len(mapping.valuations) {
			return false
		}
		valuationsCounter = 0
	}

	for symbol, value := range mapping.valuations {
		// If a variable for the symbol has not been created yet
		// then the "mapping" is creater than any "others" meaning
		if _, exists := variables[symbol]; !exists {
			return false
		}

		conjoin(symbol, value)
	}

	for _, conjunction := range valuations {
		solver.Assert(conjunction)
	}

	return solver.HasSolution()
}

func (mapping *ValuationsMap) All(
	yield func(symbol symbols.Symbol, value Value) bool,
) bool {
	for symbol, value := range mapping.valuations {
		if !yield(symbol, value) {
			return false
		}
	}
	return true
}
