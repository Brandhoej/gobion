package symbolic

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

type StatementInterpreter struct {
	expressions *ExpressionInterpreter
	sorts       *SortInterpreter
	context     *z3.Context
}

func NewStatementInterpreter(context *z3.Context) *StatementInterpreter {
	return &StatementInterpreter{
		expressions: &ExpressionInterpreter{
			context: context,
		},
		sorts: &SortInterpreter{
			context: context,
		},
		context: context,
	}
}

func (interpreter *StatementInterpreter) Fields(path *GoPath, fields *ast.FieldList) *GoPath {
	for _, field := range fields.List {
		for _, name := range field.Names {
			sort := interpreter.sorts.Expression(field.Type)
			variable := interpreter.context.NewConstant(
				z3.WithName(name.Name), sort,
			)
			path.scope.Declare(name.Name, variable, variable)
		}
	}

	return path
}

func (interpreter *StatementInterpreter) Function(path *GoPath, function *ast.FuncDecl) *GoPath {
	// Add all formal output outputs to the context.
	path = interpreter.Fields(path, function.Type.Results)
	// Add all formal input parameters to the context.
	path = interpreter.Fields(path, function.Type.Params)

	path = interpreter.Block(path, function.Body)

	return path
}

func (interpreter *StatementInterpreter) Statement(path *GoPath, statement ast.Stmt) *GoPath {
	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		return interpreter.Block(path, cast)
	case *ast.IfStmt:
		return interpreter.IfBranch(path, cast)
	/*case *ast.ReturnStmt:
	interpreter.Returns(pc, cast)*/
	case *ast.ForStmt:
		return interpreter.ForLoop(path, cast)
	case *ast.AssignStmt:
		return interpreter.Assignment(path, cast)
	case *ast.GenDecl:
	default:
		panic("Unsupported")
	}

	return path
}

func (interpreter *StatementInterpreter) Block(path *GoPath, block *ast.BlockStmt) *GoPath {
	for _, statement := range block.List {
		interpreter.Statement(path, statement)
	}

	return path
}

func (interpreter *StatementInterpreter) Assignment(path *GoPath, assignment *ast.AssignStmt) *GoPath {
	for idx := range assignment.Lhs {
		value := interpreter.expressions.Expression(path.scope, assignment.Rhs[idx]).Simplify()
		if assignment.Tok == token.ASSIGN {
			variable := interpreter.expressions.Expression(path.scope, assignment.Lhs[idx])
			path.scope.AssignTo(variable.String(), value)
		} else if assignment.Tok == token.DEFINE {
			identifier := assignment.Lhs[idx].(*ast.Ident).Name
			sort := interpreter.sorts.Expression(assignment.Rhs[idx])
			variable := path.context.NewConstant(
				z3.WithName(identifier), sort,
			)
			path.scope.Declare(identifier, variable, value)
		}
	}

	return path
}

func (interpreter *StatementInterpreter) IfBranch(path *GoPath, branch *ast.IfStmt) *GoPath {
	// Create scope from which the initialisation is available.
	if branch.Init != nil {
		interpreter.Statement(path, branch.Init)
	}

	condition := interpreter.expressions.Expression(path.scope, branch.Cond).Simplify()
	consequence := path.Branch(condition)
	if consequence != nil {
		interpreter.Block(consequence, branch.Body)
	}

	// Alternative branch (else).
	if branch.Else != nil {
		alternative := path.Branch(z3.Not(condition))
		if alternative != nil {
			interpreter.Statement(alternative, branch.Else)
		}

		// If the consequence as unsatisfiable in an if-then-else.
		// Then, we reduce it to if-else which is a if-not(then).
		if consequence == nil {
			alternative.MergeIT()
		} else if consequence != nil && alternative != nil {
			// If both branches in the if-then-else are satisfiable.
			// Then. we branch with both branches.
			return consequence.MergeITE(alternative)
		}
	}

	// We only reach this if there were no alternative branch.
	// So if the consequence was unsatisable we ignore the if-then branchin completely.
	// Otherwise, we branch normally with the consequence.
	if consequence == nil {
		return path
	} else {
		return consequence.MergeIT()
	}
}

func (interpreter *StatementInterpreter) ForLoop(path *GoPath, loop *ast.ForStmt) *GoPath {
	if loop.Init != nil {
		path = interpreter.Statement(path, loop.Init)
	}

	// The default loop condition is true. Otherwise, we interpret the loop condition and assert it.
	condition := interpreter.context.NewTrue()
	if loop.Cond != nil {
		condition = interpreter.expressions.Expression(path.scope, loop.Cond).Simplify()
	}
	if body := path.Branch(condition); body != nil {
		// After the loop condition is interpreted we interpret the loop body.
		interpreter.Block(body, loop.Body)
	
		if loop.Post != nil {
			interpreter.Statement(body, loop.Post)
		}
	
		return body.MergeIT()
	}

	return path
}