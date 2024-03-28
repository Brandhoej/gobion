package golang

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

type GoStatementInterpreter struct {
	context     *z3.Context
	sorts       *GoSortInterpreter
	expressions *GoExpressionInterpreter
	outputs     []*z3.AST
}

func NewGoStatementInterpreter(context *z3.Context) *GoStatementInterpreter {
	return &GoStatementInterpreter{
		context: context,
		sorts: &GoSortInterpreter{
			context: context,
		},
		expressions: &GoExpressionInterpreter{
			context: context,
		},
	}
}

func (interpreter *GoStatementInterpreter) statement(path *GoPath, statement ast.Stmt) *GoPath {
	switch cast := any(statement).(type) {
	case *ast.AssignStmt:
		return interpreter.assignment(path, cast)
	case *ast.IncDecStmt:
		return interpreter.incrementDecrement(path, cast)
	case *ast.ReturnStmt:
		return interpreter.returnTermination(path, cast)
	}
	panic("Unsupported statement")
}

func (interpreter *GoStatementInterpreter) incrementDecrement(path *GoPath, incDec *ast.IncDecStmt) *GoPath {
	one := interpreter.context.NewInt(1, interpreter.context.IntegerSort())
	identifier := incDec.X.(*ast.Ident).Name
	valuation, _ := path.scope.Valuation(identifier)

	if incDec.Tok == token.INC {
		path.scope.Assign(identifier, z3.Add(valuation, one))
	} else {
		path.scope.Assign(identifier, z3.Subtract(valuation, one))
	}

	return path
}

func (interpreter *GoStatementInterpreter) assignment(path *GoPath, assignment *ast.AssignStmt) *GoPath {
	for idx := range assignment.Lhs {
		identifier := assignment.Lhs[idx].(*ast.Ident).Name
		value := interpreter.expressions.Expression(path.scope, assignment.Rhs[idx]).Simplify()
		switch assignment.Tok {
		case token.ASSIGN: // =
			path.scope.Assign(identifier, value)
		case token.DEFINE: // :=
			path.scope.Declare(identifier, value)
		case token.ADD_ASSIGN: // +=
			valuation, _ := path.scope.Valuation(identifier)
			path.scope.Assign(identifier, z3.Add(valuation, value))
		case token.SUB_ASSIGN: // -=
			valuation, _ := path.scope.Valuation(identifier)
			path.scope.Assign(identifier, z3.Subtract(valuation, value))
		case token.MUL_ASSIGN: //*=
			valuation, _ := path.scope.Valuation(identifier)
			path.scope.Assign(identifier, z3.Multiply(valuation, value))
		case token.QUO_ASSIGN: // /=
			valuation, _ := path.scope.Valuation(identifier)
			path.scope.Assign(identifier, z3.Divide(valuation, value))
		case token.REM_ASSIGN: // %=
			valuation, _ := path.scope.Valuation(identifier)
			path.scope.Assign(identifier, z3.Remaninder(valuation, value))
		}
	}

	return path
}

func (interpreter *GoStatementInterpreter) returnTermination(path *GoPath, returnStatement *ast.ReturnStmt) *GoPath {
	for idx := range returnStatement.Results {
		valuation := interpreter.expressions.Expression(
			path.scope, returnStatement.Results[idx],
		)

		// If we encounter a return statement with a tautologhical PC. Then that is return value of all possible paths.
		// Otherwise, the program has atleast one branch and therefore the return value is a result of some constraints.
		// In the cases where we have multiple returns in seperate branches then the output is a if-then-else.
		// More formally but still informal: "if pc then return valuation else return existing output".
		if interpreter.outputs[idx] == nil || path.IsTautologhy() {
			interpreter.outputs[idx] = valuation.Simplify()
		} else {
			interpreter.outputs[idx] = z3.ITE(
				path.pc, valuation, interpreter.outputs[idx],
			).Simplify()
		}
	}

	return path
}
