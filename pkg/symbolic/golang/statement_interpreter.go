package golang

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbolic"
)

type GoStatementInterpreter struct {
	context     *z3.Context
	sorts       *GoSortInterpreter
	expressions *GoExpressionInterpreter
}

func NewStatementInterpreter(context *z3.Context) *GoStatementInterpreter {
	return &GoStatementInterpreter{
		context:     context,
		sorts:       NewSortInterpreter(context),
		expressions: NewExpressionsInterpreter(context),
	}
}

func (interpreter *GoStatementInterpreter) Interpret(
	scope symbolic.Scope, statements []ast.Stmt,
) (outputs []*z3.AST) {
	for idx := range statements {
		outputs = interpreter.statement(
			scope, statements[idx],
		)
	}
	return outputs
}

func (interpreter *GoStatementInterpreter) statement(
	scope symbolic.Scope, statement ast.Stmt,
) (outputs []*z3.AST) {
	switch cast := any(statement).(type) {
	case *ast.AssignStmt:
		return interpreter.assignment(scope, cast)
	case *ast.IncDecStmt:
		return interpreter.incrementDecrement(scope, cast)
	case *ast.ReturnStmt:
		return interpreter.returnTermination(scope, cast)
	}
	panic("Unsupported statement")
}

func (interpreter *GoStatementInterpreter) incrementDecrement(
	scope symbolic.Scope, incDec *ast.IncDecStmt,
) (outputs []*z3.AST) {
	one := interpreter.context.NewInt(1, interpreter.context.IntegerSort())
	identifier := incDec.X.(*ast.Ident).Name
	symbol, _ := scope.Lookup(identifier)
	valuation := scope.Valuation(symbol)

	if incDec.Tok == token.INC {
		scope.Assign(symbol, z3.Add(valuation, one))
	} else {
		scope.Assign(symbol, z3.Subtract(valuation, one))
	}

	return outputs
}

func (interpreter *GoStatementInterpreter) assignment(
	scope symbolic.Scope, assignment *ast.AssignStmt,
) (outputs []*z3.AST) {
	for idx := range assignment.Lhs {
		identifier := assignment.Lhs[idx].(*ast.Ident).Name
		value := interpreter.expressions.Expression(scope, true, assignment.Rhs[idx]).Simplify()
		switch assignment.Tok {
		case token.ASSIGN: // =
			symbol, _ := scope.Lookup(identifier)
			scope.Assign(symbol, value)
		case token.DEFINE: // :=
			symbol := scope.Bind(identifier)
			scope.Declare(symbol, value)
		case token.ADD_ASSIGN: // +=
			symbol, _ := scope.Lookup(identifier)
			valuation := scope.Valuation(symbol)
			scope.Assign(symbol, z3.Add(valuation, value))
		case token.SUB_ASSIGN: // -=
			symbol, _ := scope.Lookup(identifier)
			valuation := scope.Valuation(symbol)
			scope.Assign(symbol, z3.Subtract(valuation, value))
		case token.MUL_ASSIGN: //*=
			symbol, _ := scope.Lookup(identifier)
			valuation := scope.Valuation(symbol)
			scope.Assign(symbol, z3.Multiply(valuation, value))
		case token.QUO_ASSIGN: // /=
			symbol, _ := scope.Lookup(identifier)
			valuation := scope.Valuation(symbol)
			scope.Assign(symbol, z3.Divide(valuation, value))
		case token.REM_ASSIGN: // %=
			symbol, _ := scope.Lookup(identifier)
			valuation := scope.Valuation(symbol)
			scope.Assign(symbol, z3.Remaninder(valuation, value))
		}
	}

	return outputs
}

func (interpreter *GoStatementInterpreter) returnTermination(
	scope symbolic.Scope, returnStatement *ast.ReturnStmt,
) (outputs []*z3.AST) {
	outputs = make([]*z3.AST, len(returnStatement.Results))
	for idx := range returnStatement.Results {
		valuation := interpreter.expressions.Expression(
			scope, true, returnStatement.Results[idx],
		)
		outputs[idx] = valuation
	}

	return outputs
}
