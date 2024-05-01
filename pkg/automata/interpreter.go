package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
)

type Interpreter struct {
	context   *z3.Context
	variables language.Variables
}

func NewInterpreter(context *z3.Context, variables language.Variables) *Interpreter {
	return &Interpreter{
		context:   context,
		variables: variables,
	}
}

func (solver *Interpreter) Interpret(
	valuations language.Valuations,
	expression language.Expression,
) {
	interpreter := language.NewSymbolicInterpreter(
		solver.context, solver.variables, valuations,
	)
	interpreter.Expression(expression)
}

func (solver *Interpreter) IsSatisfied(
	valuations language.Valuations,
	expression language.Expression,
) bool {
	interpreter := language.NewSymbolicInterpreter(
		solver.context, solver.variables, valuations.Copy(),
	)
	return interpreter.Satisfies(expression)
}

func (interpreter *Interpreter) IsSatisfiable(
	expression language.Expression,
) bool {
	translator := language.NewZ3Translator(
		interpreter.context, interpreter.variables,
	)
	translation := translator.Translate(expression)
	solver := interpreter.context.NewSolver()
	return solver.HasSolutionFor(translation)
}
