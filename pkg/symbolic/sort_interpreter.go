package symbolic

import (
	"go/ast"

	"github.com/Brandhoej/gobion/internal/z3"
)

type SortInterpreter struct {
	context *z3.Context
}

func (interpreter *SortInterpreter) Expression(expression ast.Expr) *z3.Sort {
	switch cast := any(expression).(type) {
	case *ast.Ident:
		return interpreter.Identifier(cast)
	default:
		panic("Unsupported expression")
	}
}

func (interpreter *SortInterpreter) Identifier(identifier *ast.Ident) *z3.Sort {
	switch identifier.Name {
	case "bool":
		return interpreter.context.BooleanSort()
	case "int":
		return interpreter.context.IntegerSort()
	case "float":
		return interpreter.context.RealSort()
	case "false":
		return interpreter.context.BooleanSort()
	case "true":
		return interpreter.context.BooleanSort()
	default:
		panic("Unknown sort")
	}
}