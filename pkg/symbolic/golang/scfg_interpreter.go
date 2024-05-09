package golang

/*import (
	"go/ast"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/scfg/scfg"
	"github.com/Brandhoej/gobion/pkg/symbolic"
)

type scfgInterpreter struct {
	context     *z3.Context
	statements  *GoStatementInterpreter
	paths      map[int]*symbolic.Path
}

func InterpretSCFG(
	context *z3.Context,
	path *symbolic.Path,
	scopes *scfg.Graph[ast.Stmt, ast.Expr],
	cardinality int,
) (*symbolic.Path, []*z3.AST) {
	entry := scopes.CFG().Entry()
	interpreter := &scfgInterpreter{
		context: context,
		statements: NewStatementInterpreter(
			context, cardinality,
		),
		paths: map[int]*symbolic.Path{
			entry: path,
		},
	}
	return interpreter.Block(path, scopes, entry), interpreter.statements.outputs
}

func (interpreter *scfgInterpreter) Join(
	path *symbolic.Path, block int,
) *symbolic.Path {
	if existing, exists := interpreter.paths[block]; exists {
		existing.Join(path)
	} else {
		interpreter.paths[block] = path
	}

	return interpreter.paths[block]
}

func (interpreter *scfgInterpreter) Fork(source, destination int, pc *z3.AST) *symbolic.Path {
	if path, exists := interpreter.paths[destination]; exists && path != nil {
		if !path.IsFeasible() {
			return nil
		}
		return path
	}

	path := interpreter.paths[source].Fork(pc)
	interpreter.paths[destination] = path
	return path
}

func (interpreter *scfgInterpreter) Block(
	path *symbolic.Path, scopes *scfg.Graph[ast.Stmt, ast.Expr], block int,
) *symbolic.Path {
	cfg := scopes.CFG()

	if block == cfg.Exit() {
		return path
	}

	statements, condition := cfg.Block(block)
	for _, statement := range statements {
		path = interpreter.statements.statement(path, statement)
	}

	path = interpreter.Join(path, block)

	if !path.IsFeasible() {
		return path
	}

	check, jumps := cfg.Condition(condition)
	actual := interpreter.context.NewTrue()
	if check != nil {
		actual = interpreter.statements.expressions.Expression(path.Scope(), false, check)
	}

	// Now that it has been reversed it fails because the "return" is executed before the body.
	// Beucase of this it is not enabled and only after executing the body it is.
	// slices.Reverse(jumps)

	for _, jump := range jumps {
		branch, destination := cfg.Jump(jump)
		expected := interpreter.context.NewTrue()
		if branch != nil {
			expected = interpreter.statements.expressions.Expression(path.Scope(), false, branch)
		}

		equality := z3.Eq(actual, expected)
		if fork := interpreter.Fork(block, destination, equality); fork != nil {
			interpreter.Block(fork, scopes, destination)
		}
	}

	return path
}*/
