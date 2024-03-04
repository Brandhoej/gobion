package symbolic

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

type functionInterpreter struct {
	context *z3.Context
	sorts   *SortInterpreter
	expressions *GoExpressionInterpreter
	inputs []*z3.AST
}

func InterpretFunction(context *z3.Context, function *ast.FuncDecl) *z3.AST {
	if len(function.Type.Results.List) != 1 {
		panic("Declaration interpreter only supports functions with one return value.")
	}

	if len(function.Type.Results.List[0].Names) != 0 {
		panic("Declaration interpreter does not support named outputs from functions")
	}

	interpreter := functionInterpreter{
		context: context,
		sorts: &SortInterpreter{
			context: context,
		},
		expressions: &GoExpressionInterpreter{
			context: context,
		},
		inputs: make([]*z3.AST, 0),
	}

	path := NewGoGlobalPath(context)

	var input uint = 0
	for _, parameter := range function.Type.Params.List {
		sort := interpreter.sorts.Expression(parameter.Type)
		for range parameter.Names {
			variable := path.scope.Bind(input, sort)
			interpreter.inputs = append(interpreter.inputs, variable)
			input += 1
		}
	}

	path.scope.Define("ret", context.IntegerSort())

	path = interpreter.block(path, function.Body)

	returnSymbol, _ := path.scope.Symbol("ret")
	returnValuation,  _ := path.scope.valuations.Load(returnSymbol)
	return returnValuation
}

func (interpreter *functionInterpreter) statement(path *GoPath, statement ast.Stmt) *GoPath {
	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		return interpreter.block(path, cast)
	case *ast.AssignStmt:
		return interpreter.Assignment(path, cast)
	}
	panic("Unsupported statement")
}

func (interpreter *functionInterpreter) block(path *GoPath, block *ast.BlockStmt) *GoPath {
	for _, statement := range block.List {
		path = interpreter.statement(path, statement)
	}
	return path
}

func (interpreter *functionInterpreter) Assignment(path *GoPath, assignment *ast.AssignStmt) *GoPath {
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
		}
	}

	return path
}
