package symbolic

import (
	"go/ast"

	"github.com/Brandhoej/gobion/internal/z3"
)

type ApplicationInterpreter struct {
	context *z3.Context
}

func NewApplicationInterpreter(context *z3.Context) *ApplicationInterpreter {
	return &ApplicationInterpreter{
		context: context,
	}
}

func (interpreter *ApplicationInterpreter) ApplicationOf(funtion *z3.FunctionDeclaration) *z3.AST {
	return nil
}

func (interpreter *ApplicationInterpreter) statement(path *GoPath, statement ast.Stmt) *z3.AST {
	return nil
}