package symbolic

import (
	"slices"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/scfg/cfg"
	"github.com/Brandhoej/gobion/pkg/scfg/scfg"
	"github.com/Brandhoej/gobion/pkg/structures"
)

type StatementInterpreter[T any] interface {
	Interpret(scope Scope, statements []T) (outputs []*z3.AST)
}

type ExpressionInterpreter[T any] interface {
	Interpret(scope Scope, expression T, reduced bool) *z3.AST
}

type CallInterpreter interface {
	Interpret(path *Path) (*Path, []*z3.AST)
}

type SCFG[S any, E comparable] struct {
	context *z3.Context
	graph scfg.Graph[S, E]
	statements StatementInterpreter[S]
	expressions ExpressionInterpreter[E]
}

type work struct {
	path *Path
	block int
}

func NewSCFG[S any, E comparable](
	context *z3.Context,
	graph scfg.Graph[S, E],
	statements StatementInterpreter[S],
	expressions ExpressionInterpreter[E],
) *SCFG[S, E] {
	return &SCFG[S, E]{
		context: context,
		graph: graph,
		statements: statements,
		expressions: expressions,
	}
}

func (scfg *SCFG[S, E]) Guard(
	scope Scope, condition, jump int, cfg *cfg.Graph[S, E],
) *z3.AST {
	check, _ := cfg.Condition(condition)
	actual := scfg.expressions.Interpret(scope, check, false)
	if actual == nil {
		actual = scfg.context.NewTrue()
	}

	branch, _ := cfg.Jump(jump)
	expected := scfg.expressions.Interpret(scope, branch, false)
	if expected == nil {
		expected = scfg.context.NewTrue()
	}

	return z3.Eq(actual, expected)
}

func (scfg *SCFG[S, E]) Interpret(path *Path) (*Path, []*z3.AST) {
	outputs := make([]*z3.AST, 0)
	scopes := scfg.graph
	cfg := scopes.CFG()

	worklist := make(structures.Queue[*work], 0)
	worklist.Enqueue(&work{
		path: path,
		block: cfg.Entry(),
	})

	next := func() *work {
		var idx int = -1
		for i, work := range worklist {
			if work.path.IsFeasible() {
				idx = i
				break
			}
		}

		if idx == -1 {
			return nil
		}

		work := worklist[idx]
		worklist = slices.Delete(worklist, idx, idx+1)
		return work
	}

	output := func(pc *z3.AST, returns []*z3.AST) {
		for idx := range returns {
			if len(outputs) == 0 {
				outputs = returns
			} else {
				outputs[idx] = z3.ITE(
					pc, returns[idx], outputs[idx],
				)
			}
		}
	}

	for !worklist.IsEmpty() {
		// Step 1: Get the first front that we can execute.
		front := next()
		if front == nil {
			break
		}
		if front.block == cfg.Exit() {
			continue
		}

		// Step 2: Execute the front. We know it is feasible. Combine outputs.
		statements, condition := cfg.Block(front.block)
		returns := scfg.statements.Interpret(front.path.scope, statements)
		output(front.path.pc, returns)

		// Step 3: Add all jumps to the worklist.
		_, jumps := cfg.Condition(condition)
		for _, jump := range jumps {
			guard := scfg.Guard(
				front.path.scope, condition, jump, cfg,
			)
			fork := front.path.Fork(guard)
			_, destination := cfg.Jump(jump)
			worklist.Enqueue(&work{
				path: fork,
				block: destination,
			})
		}
	}

	return path, outputs
}