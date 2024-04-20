package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type ConstraintSolver struct {
	constraints constraints.SymbolicInterpreter
	expressions expressions.SymbolicInterpreter
	valuations  expressions.Valuations
	variables   expressions.Variables
	backing     *z3.Solver
}

func NewConstraintSolver(
	solver *z3.Solver, valuations expressions.Valuations, variables expressions.Variables,
) *ConstraintSolver {
	return &ConstraintSolver{
		constraints: constraints.NewSymbolicInterpreter(
			solver.Context(), valuations,
		),
		expressions: expressions.NewSymbolicInterpreter(
			solver.Context(), valuations,
		),
		valuations: valuations,
		variables: variables,
		backing: solver,
	}
}

func (solver *ConstraintSolver) setup() {
	solver.valuations.All(func(symbol symbols.Symbol, value expressions.Expression) bool {
		if sort, exists := solver.variables.Lookup(symbol); exists {
			context := solver.backing.Context()
			var z3Sort *z3.Sort
			switch sort {
			case expressions.IntegerSort:
				z3Sort = context.IntegerSort()
			case expressions.BooleanSort:
				z3Sort = context.BooleanSort()
			}
			constant := context.NewConstant(
				z3.WithInt(int(symbol)), z3Sort,
			)
			valuation := solver.expressions.Interpret(value)
			solver.backing.Assert(z3.Eq(constant, valuation))
		} else {
			panic("Variable with symbol not declared")
		}

		return true
	})
}

func (solver *ConstraintSolver) HasSolutionFor(constraint constraints.Constraint) bool {
	solver.backing.Push()
	defer solver.backing.Pop(1)
	solver.setup()

	proposition := solver.constraints.Interpret(constraint)
	return solver.backing.HasSolutionFor(proposition)
}