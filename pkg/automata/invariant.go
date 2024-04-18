package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Invariant struct {
	context   *z3.Context
	condition _Expression
}

func NewInvariant(context *z3.Context, condition _Expression) Invariant {
	return Invariant{
		context:   context,
		condition: condition,
	}
}

func (invariant Invariant) IsSatisfiable(variables Variables, valuations Valuations) bool {
	solver := newSolver(invariant.context.NewSolver())
	solver.Assert(invariant.condition)

	valuations.All(func(symbol symbols.Symbol, value Value) bool {
		if variable, exists := variables.Variable(symbol); exists {
			expression := newExpression(value.relation.ast(variable, value.ast))
			solver.Assert(expression)
		} else {
			panic("Unknown variable symbol")
		}

		return true
	})

	return solver.HasSolution()
}
