package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Solver struct {
	context *z3.Context
	backing *z3.Solver
	variables language.Variables
	valuations language.Valuations
	expressions language.SymbolicExpressionInterpreter
}

func newSolver(
	solver *z3.Solver,
	variables language.Variables,
	valuations language.Valuations,
) Solver {
	return Solver{
		context: solver.Context(),
		backing: solver,
		variables: variables,
		valuations: valuations,
		expressions: language.NewSymbolExpressionInterpreter(
			solver.Context(), valuations,
		),
	}
}

func (solver Solver) setup() {
	solver.valuations.All(func(symbol symbols.Symbol, value language.Expression) bool {
		if variable, exists := solver.variables.Variable(symbol); exists {
			solver.Assert(language.NewBinary(variable, language.Equal, value))
		} else {
			panic("Unknown variable")
		}
		return true
	})
}

func (solver Solver) Assert(expression language.Expression) {
	proposition := solver.expressions.Interpret(expression)
	solver.backing.Assert(proposition)
}

func (solver Solver) HasSolutionFor(expression language.Expression) bool {
	solver.backing.Push()
	defer solver.backing.Pop(1)
	solver.setup()

	proposition := solver.expressions.Interpret(expression)
	return solver.backing.HasSolutionFor(proposition)
}

func (solver Solver) HasSolution() bool {
	solver.backing.Push()
	defer solver.backing.Pop(1)
	solver.setup()

	return solver.backing.HasSolution()
}
