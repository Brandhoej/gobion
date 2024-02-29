package symbolic

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

/*
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
	expressions     *ExpressionInterpreter
	solver          *z3.Solver
	context         *z3.Context
	returnVariables []*z3.AST
	awaited         []func()
	visited         func(node ast.Node)
}

func WithSolver(solver *z3.Solver) functionInterpreterConfig {
	return func(interpreter *FunctionInterpreter) {
		interpreter.solver = solver
	}
}

func Visited(observer func(node ast.Node)) functionInterpreterConfig {
	return func(interpreter *FunctionInterpreter) {
		interpreter.visited = observer
	}
}

func NewFunctionInterpreter(context *z3.Context, function *ast.FuncDecl, configs ...functionInterpreterConfig) *FunctionInterpreter {
	// Construct the return variables.
	parameters := function.Type.Params
	returns := make([]*z3.AST, 0, len(parameters.List))
	for outer, field := range parameters.List {
		for inner, name := range field.Names {
			returns[outer+inner] = context.NewConstant(
				z3.WithName(name.Name), context.BooleanSort(),
			)
		}
	}

	// Create the default interpreter.
	interpreter := &FunctionInterpreter{
		expressions: &ExpressionInterpreter{
			context: context,
		},
		solver:          nil,
		context:         context,
		returnVariables: returns,
		awaited:         make([]func(), 0),
	}

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

func (interpreter *FunctionInterpreter) Step() {
	length := len(interpreter.awaited)
	if length == 0 {
		return
	}

	// We want a stack like behvaiour so we take from the end and remove the last element.
	interpreter.awaited[length-1]()
	interpreter.awaited = interpreter.awaited[0 : length-1]
}

func (interpreter *FunctionInterpreter) await(awaited func()) {
	interpreter.awaited = append(interpreter.awaited, awaited)
}

func (interpreter *FunctionInterpreter) function(function *ast.FuncDecl) {
	interpreter.await(func() {
		interpreter.solver.Push()
		interpreter.block(function.Body)
		interpreter.visited(function.Body)
		interpreter.solver.Pop(1)
	})
}

func (interpreter *FunctionInterpreter) statement(statement ast.Stmt) {
	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		interpreter.block(cast)
	case *ast.IfStmt:
		interpreter.ifBranch(cast)
	case *ast.ReturnStmt:
		interpreter.returns(cast)
	case *ast.ForStmt:
		interpreter.forLoop(cast)
	}
	panic("Unsupported")
}

func (interpreter *FunctionInterpreter) block(block *ast.BlockStmt) {
	interpreter.await(func() {
		interpreter.solver.Push()
		for _, statement := range block.List {
			interpreter.statement(statement)
		}
		interpreter.visited(block)
		interpreter.solver.Pop(1)
	})
}

func (interpreter *FunctionInterpreter) ifBranch(branch *ast.IfStmt) {
	hasAlternative := branch.Else != nil

	// Create scope from which the initialisation is available.
	if branch.Init != nil {
		interpreter.await(func() {
			interpreter.solver.Push()
			interpreter.statement(branch.Init)
			interpreter.visited(branch.Init)
		})
	}

	var condition *z3.AST
	interpreter.await(func() {
		// We set the condition in the consequence such that the alternative can use it.
		// We dont need a seperate awaited execution for the condition.
		// The reason for this is that the condition will always be present.
		condition = interpreter.expressions.Expression(branch.Cond)

		// Create scope of the consequence which is not shared between this and the optional alternative branch.
		interpreter.solver.Push()
		interpreter.solver.Assert(condition)
		interpreter.block(branch.Body)
		interpreter.visited(branch.Body)
		interpreter.solver.Pop(1)

		// Pop the initilisation scope in the consequence if we do not have an alternative branch.
		if !hasAlternative {
			interpreter.solver.Pop(1)
		}
	})

	// Alternative branch (else).
	if hasAlternative {
		interpreter.await(func() {
			interpreter.solver.Push()
			interpreter.solver.Assert(z3.Not(condition))
			interpreter.statement(branch.Else)
			interpreter.visited(branch.Else)
			interpreter.solver.Pop(1)

			// Pop the initilisation scope.
			interpreter.solver.Pop(1)
		})
	}
}

func (interpreter *FunctionInterpreter) forLoop(loop *ast.ForStmt) {
	hasInitialisation := loop.Init != nil
	hasConditional := loop.Cond != nil
	hasUpdate := loop.Post != nil

	if hasInitialisation {
		interpreter.await(func() {
			// Create for loop scope with potentially new variables from initialisation.
			interpreter.solver.Push()
			interpreter.statement(loop.Init)
			interpreter.visited(loop.Init)
		})
	}

	interpreter.await(func() {
		// If we did not have an initialisation we need to create the scope when to interpret the body.
		if !hasInitialisation {
			interpreter.solver.Push()
		}

		// The default loop condition is true. Otherwise, we interpret the loop condition and assert it.
		condition := interpreter.context.NewTrue()
		if hasConditional {
			condition = interpreter.expressions.Expression(loop.Cond)
		}
		interpreter.solver.Assert(condition)

		// After the loop condition is interpreted we interpret the loop body.
		interpreter.block(loop.Body)
		interpreter.visited(loop.Body)

		// If we dont have an update we have to close the scope after the body.
		if !hasUpdate {
			interpreter.solver.Pop(1)
		}
	})

	interpreter.await(func() {
		if hasUpdate {
			interpreter.statement(loop.Post)
			interpreter.visited(loop.Post)
			
			// The update has to close the scope after it has been interpreted. Otherwise, the body closes it.
			interpreter.solver.Pop(1)
		}
	})
}

func (interpreter *FunctionInterpreter) returns(exit *ast.ReturnStmt) {
	interpreter.await(func() {
		for idx, result := range exit.Results {
			expr := interpreter.expressions.Expression(result)
			returnVar := interpreter.returnVariables[idx]
			interpreter.solver.Assert(z3.Eq(returnVar, expr))
		}
		interpreter.visited(exit)
	})
}

type ExpressionInterpreter struct {
	context *z3.Context
}

func (interpreter *ExpressionInterpreter) Expression(expression ast.Expr) *z3.AST {
	switch cast := any(expression).(type) {
	case *ast.BinaryExpr:
		return interpreter.Binary(cast)
	case *ast.UnaryExpr:
		return interpreter.Unary(cast)
	case *ast.ParenExpr:
		return interpreter.Parenthesis(cast)
	}
	panic("Unsupported")
}

func (interpreter *ExpressionInterpreter) Binary(binary *ast.BinaryExpr) *z3.AST {
	lhs, rhs := interpreter.Expression(binary.X), interpreter.Expression(binary.Y)
	switch binary.Op {
	case token.AND:
		return z3.And(lhs, rhs)
	case token.OR:
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
	}
	panic("Unsupported")
}
