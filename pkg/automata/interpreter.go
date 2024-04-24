package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
)

type Interpreter struct {
	context   *z3.Context
	variables expressions.Variables[*z3.Sort]
}

func NewInterpreter(context *z3.Context, variables expressions.Variables[*z3.Sort]) *Interpreter {
	return &Interpreter{
		context:   context,
		variables: variables,
	}
}

func (solver *Interpreter) Interpret(
	valuations expressions.Valuations[*z3.AST],
	expression expressions.Expression,
) {
	interpreter := expressions.NewSymbolicInterpreter(
		solver.context, solver.variables, valuations,
	)
	interpreter.Interpret(expression)
}

func (solver *Interpreter) IsSatisfied(
	valuations expressions.Valuations[*z3.AST],
	expression expressions.Expression,
) bool {
	interpreter := expressions.NewSymbolicInterpreter(
		solver.context, solver.variables, valuations.Copy(),
	)
	return interpreter.Satisfies(expression)
}

func (interpreter *Interpreter) IsSatisfiable(
	expression expressions.Expression,
) bool {
	translator := expressions.NewZ3Translator(
		interpreter.context, interpreter.variables,
	)
	translation := translator.Translate(expression)
	solver := interpreter.context.NewSolver()
	return solver.HasSolutionFor(translation)
}
