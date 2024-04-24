package statements

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/automata/language/state"
)

type SymbolicInterpreter struct {
	valuations state.Valuations[*z3.AST]
	expressions expressions.SymbolicInterpreter
}

func NewSymbolicInterpreter(
	context *z3.Context,
	variables state.Variables[*z3.Sort],
	valuations state.Valuations[*z3.AST],
) *SymbolicInterpreter {
	return &SymbolicInterpreter{
		valuations: valuations,
		expressions: *expressions.NewSymbolicInterpreter(
			context, variables, valuations,
		),
	}
}

func (interpreter *SymbolicInterpreter) Interpret(statement Statement) {
	switch cast := any(statement).(type) {
	case Assignment:
		interpreter.Assignment(cast)
	}
}

func (interpreter *SymbolicInterpreter) Assignment(assignment Assignment) {
	valuation := interpreter.expressions.Interpret(assignment.valuation)
	symbol := assignment.variable.Symbol()
	interpreter.valuations.Assign(symbol, valuation)
}