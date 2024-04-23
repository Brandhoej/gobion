package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type ConstraintSolver struct {
	backing     *z3.Solver
	variables   expressions.Variables[*z3.Sort]
}

func NewConstraintSolver(solver *z3.Solver, variables expressions.Variables[*z3.Sort]) *ConstraintSolver {
	return &ConstraintSolver{
		backing: solver,
		variables: variables,
	}
}

func (solver *ConstraintSolver) setup(valuations expressions.Valuations[*z3.AST]) {
	context := solver.backing.Context()
	valuations.All(func(symbol symbols.Symbol, value *z3.AST) bool {
		if sort, exists := solver.variables.Lookup(symbol); exists {
			constant := context.NewConstant(z3.WithInt(int(symbol)), sort)
			equality := z3.Eq(constant, value)
			solver.backing.Assert(equality)
		}

		return true
	})
}

func (solver *ConstraintSolver) Satisfies(
	valuations expressions.Valuations[*z3.AST], constraint constraints.Constraint,
) bool {
	solver.backing.Push()
	defer solver.backing.Pop(1)
	solver.setup(valuations)

	interpreter := constraints.NewSymbolicInterpreter(
		solver.backing.Context(), valuations, solver.variables,
	)
	proposition := interpreter.Interpret(constraint)
	solver.backing.Assert(proposition)
	return solver.backing.HasSolutionFor(proposition)
}
