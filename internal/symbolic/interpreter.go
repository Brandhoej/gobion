package symbolic

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

/*
TODO:
- Convert a function to a symbolic representation.

func Foo(a, b bool) bool {
	if a && !b {
		return true
	}
	return b || a
}

func Moo() bool {
	if true {
		return false
	}
	return true
}

func Bar(x, y int) int {
	if x < 10 {
		return x + y
	}
	return x - y
}
*/

type functionInterpreterConfig func(interpreter *FunctionInterpreter)

type FunctionInterpreter struct {
	expressions *ExpressionInterpreter
	solver      *z3.Solver
	context     *z3.Context
	variables   map[string]*z3.AST
	returns     []string
	pc          *z3.AST
}

func WithSolver(solver *z3.Solver) functionInterpreterConfig {
	return func(interpreter *FunctionInterpreter) {
		interpreter.solver = solver
	}
}

func NewFunctionInterpreter(context *z3.Context, function *ast.FuncDecl, configs ...functionInterpreterConfig) *FunctionInterpreter {
	// Create the default interpreter.
	variables := make(map[string]*z3.AST)
	interpreter := &FunctionInterpreter{
		expressions: &ExpressionInterpreter{
			context:   context,
			variables: variables,
		},
		solver:    nil,
		context:   context,
		variables: variables,
	}

	// Construct the return variables.
	parameters := function.Type.Results.List
	returns := make([]string, 0)
	for _, field := range parameters {
		for _, name := range field.Names {
			variable := context.NewConstant(
				z3.WithName(name.Name), context.BooleanSort(),
			)
			interpreter.variables[name.Name] = variable
			returns = append(returns, name.Name)
		}
	}
	interpreter.returns = returns

	// Configure the interpreter with values that are optional.
	for _, config := range configs {
		config(interpreter)
	}

	// For optimization we wait until now to set the solver if a config with an existing one was not specified.
	if interpreter.solver == nil {
		interpreter.solver = context.NewSolver()
	}

	// Start the interpreter. This is awaited immeidately and requires a Step call by the user.
	interpreter.function(function)

	return interpreter
}

func (interpreter *FunctionInterpreter) function(function *ast.FuncDecl) {
	// Add all formal input parameters to the context.
	parameters := function.Type.Params.List
	for _, parameter := range parameters {
		for _, name := range parameter.Names {
			variable := interpreter.context.NewConstant(
				z3.WithName(name.Name), interpreter.context.BooleanSort(),
			)
			interpreter.variables[name.Name] = variable
		}
	}

	interpreter.Block(function.Body)
}

func (interpreter *FunctionInterpreter) Statement(statement ast.Stmt) {
	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		interpreter.Block(cast)
	case *ast.IfStmt:
		interpreter.IfBranch(cast)
	case *ast.ReturnStmt:
		interpreter.Returns(cast)
	case *ast.ForStmt:
		interpreter.ForLoop(cast)
	case *ast.AssignStmt:
		interpreter.Assignment(cast)
	default:
		panic("Unsupported")
	}
}

func (interpreter *FunctionInterpreter) Block(block *ast.BlockStmt) {
	for _, statement := range block.List {
		interpreter.Statement(statement)
	}
}

func (interpreter *FunctionInterpreter) IfBranch(branch *ast.IfStmt) {
	// Create scope from which the initialisation is available.
	if branch.Init != nil {
		interpreter.Statement(branch.Init)
	}

	var condition *z3.AST
	// We set the condition in the consequence such that the alternative can use it.
	// We dont need a seperate awaited execution for the condition.
	// The reason for this is that the condition will always be present.
	condition = interpreter.expressions.Expression(branch.Cond)

	// Create scope of the consequence which is not shared between this and the optional alternative branch.
	interpreter.solver.Assert(condition)
	interpreter.Block(branch.Body)

	interpreter.solver.Assert(z3.Not(condition))

	// Alternative branch (else).
	if branch.Else != nil {
		interpreter.Statement(branch.Else)
	}
}

func (interpreter *FunctionInterpreter) ForLoop(loop *ast.ForStmt) {
	if loop.Init != nil {
		interpreter.Statement(loop.Init)
	}

	// The default loop condition is true. Otherwise, we interpret the loop condition and assert it.
	condition := interpreter.context.NewTrue()
	if loop.Cond != nil {
		condition = interpreter.expressions.Expression(loop.Cond)
	}
	interpreter.solver.Assert(condition)

	// After the loop condition is interpreted we interpret the loop body.
	interpreter.Block(loop.Body)

	if loop.Post != nil {
		interpreter.Statement(loop.Post)
	}
}

func (interpreter *FunctionInterpreter) Returns(exit *ast.ReturnStmt) {
	if len(exit.Results) == 0 {
		return
	}

	for idx, result := range exit.Results {
		expr := interpreter.expressions.Expression(result)
		namedReturn := interpreter.returns[idx]
		returnVariable := interpreter.variables[namedReturn]
		assignment := z3.Eq(returnVariable, expr)
		interpreter.solver.Assert(assignment)
	}

	fmt.Println(interpreter.solver.String())
}

func (interpreter *FunctionInterpreter) Assignment(assignment *ast.AssignStmt) {
	if assignment.Tok == token.ASSIGN {
		for idx := range assignment.Lhs {
			lhs := interpreter.expressions.Expression(assignment.Lhs[idx])
			rhs := interpreter.expressions.Expression(assignment.Rhs[idx])
			equality := z3.Eq(lhs, rhs)
			interpreter.solver.Assert(equality)
		}
	}
}

type ExpressionInterpreter struct {
	context   *z3.Context
	variables map[string]*z3.AST
}

func (interpreter *ExpressionInterpreter) Expression(expression ast.Expr) *z3.AST {
	switch cast := any(expression).(type) {
	case *ast.BinaryExpr:
		return interpreter.Binary(cast)
	case *ast.UnaryExpr:
		return interpreter.Unary(cast)
	case *ast.ParenExpr:
		return interpreter.Parenthesis(cast)
	case *ast.Ident:
		return interpreter.Identifier(cast)
	}
	panic("Unsupported")
}

func (interpreter *ExpressionInterpreter) Binary(binary *ast.BinaryExpr) *z3.AST {
	lhs, rhs := interpreter.Expression(binary.X), interpreter.Expression(binary.Y)
	switch binary.Op {
	case token.LAND:
		return z3.And(lhs, rhs)
	case token.LOR:
		return z3.Or(lhs, rhs)
	case token.XOR:
		return z3.Xor(lhs, rhs)
	}
	panic("Unsupported")
}

func (interpreter *ExpressionInterpreter) Parenthesis(parenthesis *ast.ParenExpr) *z3.AST {
	return interpreter.Expression(parenthesis.X)
}

func (interpreter *ExpressionInterpreter) Unary(unary *ast.UnaryExpr) *z3.AST {
	operand := interpreter.Expression(unary.X)
	switch unary.Op {
	case token.NOT:
		return z3.Not(operand)
	}
	panic("Unsupported")
}

func (interpreter *ExpressionInterpreter) Identifier(identifier *ast.Ident) *z3.AST {
	switch identifier.Name {
	case "true":
		return interpreter.context.NewTrue()
	case "false":
		return interpreter.context.NewFalse()
	default:
		return interpreter.variables[identifier.Name]
	}
}
