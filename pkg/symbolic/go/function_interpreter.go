package symbolic

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

type functionInterpreter struct {
	context     *z3.Context
	sorts       *SortInterpreter
	expressions *GoExpressionInterpreter
	outputs     []*z3.AST
}

func InterpretFunction(path *GoPath, function *ast.FuncDecl, inputs []*z3.AST) []*z3.AST {
	if len(function.Type.Results.List[0].Names) != 0 {
		panic("Declaration interpreter does not support named outputs from functions")
	}

	interpreter := functionInterpreter{
		context: path.context,
		sorts: &SortInterpreter{
			context: path.context,
		},
		expressions: &GoExpressionInterpreter{
			context: path.context,
		},
		outputs: make([]*z3.AST, 1),
	}

	path = path.Enclose()

	var index uint = 0
	for _, parameter := range function.Type.Params.List {
		for _, name := range parameter.Names {
			path.scope.Declare(name.Name, inputs[index])
			index += 1
		}
	}

	interpreter.block(path, function.Body)

	return interpreter.outputs
}

func (interpreter *functionInterpreter) statement(path *GoPath, statement ast.Stmt) *GoPath {
	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		return interpreter.block(path, cast)
	case *ast.ReturnStmt:
		return interpreter.returnTermination(path, cast)
	case *ast.IfStmt:
		return interpreter.ifBranch(path, cast)
	case *ast.AssignStmt:
		return interpreter.assignment(path, cast)
	case *ast.ForStmt:
		return interpreter.forLoop(path, cast)
	case *ast.IncDecStmt:
		return interpreter.incrementDecrement(path, cast)
	}
	panic("Unsupported statement")
}

func (interpreter *functionInterpreter) block(path *GoPath, block *ast.BlockStmt) *GoPath {
	for _, statement := range block.List {
		path = interpreter.statement(path, statement)

		// A return statement was interpreted so we terminate the block interpretation.
		if path.IsTerminated() {
			return path
		}
	}
	return path
}

func (interpreter *functionInterpreter) returnTermination(path *GoPath, returnStatement *ast.ReturnStmt) *GoPath {
	for idx := range returnStatement.Results {
		valuation := interpreter.expressions.Expression(
			path.scope, returnStatement.Results[idx],
		)

		// If we encounter a return statement with a tautologhical PC. Then that is return value of all possible paths.
		// Otherwise, the program has atleast one branch and therefore the return value is a result of some constraints.
		// In the cases where we have multiple returns in seperate branches then the output is a if-then-else.
		// More formally but still informal: "if pc then return valuation else return existing output".
		if interpreter.outputs[idx] == nil || path.IsTautologhy() {
			interpreter.outputs[idx] = valuation
		} else {
			interpreter.outputs[idx] = z3.ITE(
				path.pc, valuation, interpreter.outputs[idx],
			)
		}
	}

	// Terminating the paths tell the callee that the path has returned.
	path.Terminate()

	return path
}

func (interpreter *functionInterpreter) ifBranch(path *GoPath, branch *ast.IfStmt) *GoPath {
	enclosure := path.Enclose()
	if branch.Init != nil {
		enclosure = interpreter.statement(enclosure, branch.Init)
	}

	condition := interpreter.expressions.Expression(enclosure.scope, branch.Cond).Simplify()
	consequence := enclosure.Branch(condition)
	if consequence != nil {
		consequence = interpreter.block(consequence, branch.Body)
	}

	if branch.Else != nil {
		alternative := path.Branch(z3.Not(condition))
		if alternative != nil {
			alternative = interpreter.statement(alternative, branch.Else)
		}

		// If the consequence as unsatisfiable in an if-then-else.
		// Then, we reduce it to if-else which is a if-not(then).
		if consequence == nil {
			alternative.MergeIT()
		} else if alternative != nil {
			// If both branches in the if-then-else are satisfiable.
			// Then. we branch with both branches.
			return consequence.MergeITE(alternative)
		}
	}

	// We only reach this if there were no alternative branch.
	// So if the consequence was unsatisable we ignore the if-then branching completely.
	// Otherwise, we branch normally with the consequence.
	if consequence == nil {
		return path
	} else {
		return consequence.MergeIT()
	}
}

func (interpreter *functionInterpreter) forLoop(path *GoPath, loop *ast.ForStmt) *GoPath {
	// We have to declare the initialisation in an enclosure that is the parent of the body.
	// Because we want the opportunity to shadow declared variables from the enclosure in the body.
	enclosure := path.Enclose()
	if loop.Init != nil {
		enclosure = interpreter.statement(enclosure, loop.Init)
	}

	body := enclosure.Enclose()

	// There is no gurantee for the loop to ever terminate.
	// For this reason, we want to make the unfolding of it bounded to some max value.
	// For now the max unfolding is 100. But this should be based on a heuristic.
	// TODO: Make a guess? Often for loops are of the form "for i := 0; i < 10; i++".
	// Wy make a guess? It would be sad to limit to 100 unfoldings but the 101th would have found the meaning of life.
	counter, bound := 0, 100

	// Single Static Assignment is ensure by the scope handling of assignment.
	// For this reason unfolding is trivial as it is handle inherently.
	for {
		// If we exceeds to max bounded unfolding then stop.
		if counter >= bound {
			break
		}

		// If there is a loop condition we evaluate it on the enclusre path/scope.
		// We should not evaluate it on the body as variables declared in the body
		// cannot be used in the condition.
		if loop.Cond != nil {
			condition := interpreter.expressions.Expression(enclosure.scope, loop.Cond)
			body = enclosure.Branch(condition)
		}

		// If the body is nil then the branching failed most likely because the condition made the path infeasible.
		// If the body is infeasible then then the condition might be unsatisfiable.
		if body == nil || body.IsInfeasible() {
			break
		}

		// The body is executed in the body path/scope. Which is the child of the enclosure.
		if loop.Body != nil {
			body = interpreter.block(body, loop.Body)
		}

		// The post containing the update following the body - it is often an increment or decrement.
		// The post is evaluated inside the enclosure and not the body to ensure correct Go semantics.
		if loop.Post != nil {
			enclosure = interpreter.statement(enclosure, loop.Post)
		}

		// Increment the counter to keep track of the unfolding count.
		counter++
	}

	return path
}

func (interpreter *functionInterpreter) incrementDecrement(path *GoPath, incDec *ast.IncDecStmt) *GoPath {
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

func (interpreter *functionInterpreter) assignment(path *GoPath, assignment *ast.AssignStmt) *GoPath {
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
		case token.REM_ASSIGN: // /%
			valuation, _ := path.scope.Valuation(identifier)
			path.scope.Assign(identifier, z3.Remaninder(valuation, value))
		}
	}

	return path
}
