package symbolic

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

type SortInterpreter struct {
	context *z3.Context
}

func (interpreter *SortInterpreter) Expression(expression ast.Expr) *z3.Sort {
	switch cast := any(expression).(type) {
	case *ast.Ident:
		return interpreter.Identifier(cast)
	case *ast.BasicLit:
		return interpreter.Literal(cast)
	}
	panic("Unsupported expression")
}

func (interpreter *SortInterpreter) Literal(literal *ast.BasicLit) *z3.Sort {
	switch literal.Kind {
	case token.INT:
		return interpreter.context.IntegerSort()
	}
	panic("Unsupported literal type")
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
	}
	panic("Unknown sort")
}
