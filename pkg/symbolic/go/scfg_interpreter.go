package golang

import (
	"go/ast"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/scfg/scfg"
)

type scfgInterpreter struct {
	context     *z3.Context
	sorts       *GoSortInterpreter
	expressions *GoExpressionInterpreter
	statements  *GoStatementInterpreter
	scopes      map[int]*GoScope
}

func InterpretSCFG(
	path *GoPath, scopes *scfg.Graph[ast.Stmt, ast.Expr], cardinality int,
) (*GoPath, []*z3.AST) {
	sorts := &GoSortInterpreter{
		context: path.context,
	}
	expressions := &GoExpressionInterpreter{
		context: path.context,
	}
	statements := &GoStatementInterpreter{
		context: path.context,
		outputs: make([]*z3.AST, cardinality),
		expressions: expressions,
	}
	interpreter := &scfgInterpreter{
		context: path.context,
		sorts: sorts,
		expressions: expressions,
		statements: statements,
		scopes: map[int]*GoScope{
			scopes.Global(): path.scope,
		},
	}
	entry := scopes.CFG().Entry()
	return interpreter.Block(path, scopes, entry), interpreter.statements.outputs
}

func (interpreter *scfgInterpreter) Block(
	path *GoPath, flow *scfg.Graph[ast.Stmt, ast.Expr], block int,
) *GoPath {
	if block == flow.CFG().Exit() {
		return path
	}

	statements, condition := flow.CFG().Block(block)
	for _, statement := range statements {
		path = interpreter.statements.statement(path, statement)
	}

	if path.IsFeasible() {
		return interpreter.Condition(path, flow, condition)
	} 
	return path
}

func (interpreter *scfgInterpreter) Condition(
	path *GoPath, flow *scfg.Graph[ast.Stmt, ast.Expr], condition int,
) *GoPath {
	check, jumps := flow.CFG().Condition(condition)
	actual := path.context.NewTrue()
	if check != nil {
		actual = interpreter.expressions.Expression(path.scope, check)
	}

	for _, jump := range jumps {
		path = interpreter.Jump(path, flow, actual, jump)
	}

	return path
}

func (interpreter *scfgInterpreter) Jump(
	path *GoPath, flow *scfg.Graph[ast.Stmt, ast.Expr],
	actual *z3.AST, jump int,
) *GoPath {
	branch, destination := flow.CFG().Jump(jump)
	if destination == flow.CFG().Exit() {
		return path
	}

	expected := path.context.NewTrue()
	if branch != nil {
		expected = interpreter.expressions.Expression(path.scope, branch)
	}

	// Check if the jump is enabled by equality of actual and expected.
	// That is, the condition must have a solution to the 
	// rhs: true, false, case conditions, or somethign else.
	// Example 1: actual: V[i -> 1]: actual (condition): i == 1. expected (jump): true. actual == expected (Sat).
	// Example 2: actual: V[i -> 3]: actual (condition): i == 2. expected (jump): true. actual != expected (Unsat).
	// We return immediately because we assume determinsitic programs.
	// For this reason only a single jump can be "enabled".
	scopeID, _ := flow.ScopeWith(destination)
	if branch := path.Branch(z3.Eq(actual, expected), scopeID); branch != nil {
		path = interpreter.Block(branch, flow, destination).MergeIT()
	}

	return path
}
