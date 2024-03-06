package symbolic

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
)

type GoPath struct {
	parent     *GoPath
	scope      *GoScope
	context    *z3.Context
	pc         *z3.AST
	terminated bool
}

func NewGoGlobalPath(context *z3.Context) *GoPath {
	return &GoPath{
		parent:  nil,
		scope:   NewGoGlobalScope(),
		context: context,
		pc:      context.NewTrue(),
	}
}

func (path *GoPath) Terminate() {
	path.terminated = true
}

func (path *GoPath) IsTerminated() bool {
	return path.terminated
}

func (path *GoPath) Solver() *z3.Solver {
	solver := path.context.NewSolver()
	// Provided the symbolic states we have and our path constraint.
	// Is there always a solution.
	for _, identifier := range path.scope.Identifiers() {
		variable, _ := path.scope.Variable(identifier)
		valuation, _ := path.scope.Valuation(identifier)
		solver.Assert(z3.Eq(variable, valuation))
	}
	solver.Assert(path.pc)
	return solver
}

func (path *GoPath) IsTautologhy() bool {
	solver := path.Solver()
	return solver.Prove(path.context.NewTrue()) == nil
}

func (path *GoPath) IsFeasible() bool {
	solver := path.Solver()
	return solver.HasSolution()
}

func (path *GoPath) IsInfeasible() bool {
	solver := path.Solver()
	return !solver.HasSolution()
}

func (path *GoPath) Enclose() *GoPath {
	return path.Branch(path.pc)
}

func (path *GoPath) Branch(condition *z3.AST) *GoPath {
	branch := &GoPath{
		parent:     path,
		scope:      path.scope.Child(),
		context:    path.context,
		pc:         z3.And(path.pc, condition).Simplify(),
		terminated: path.terminated,
	}

	// After updating the path constraint check that it is still feasible.
	if branch.IsInfeasible() {
		return nil
	}

	return branch
}

// Merges an if-then-else branch.
func (consequence *GoPath) MergeITE(alternative *GoPath) *GoPath {
	parent := consequence.parent

	// Q: Should there be a debug assertion that no junction is infeasible?
	junction := &GoPath{
		parent:     parent,
		scope:      parent.scope,
		context:    parent.context,
		pc:         parent.pc,
		terminated: false,
	}

	solver := consequence.context.NewSolver()

	// The junction has been terminated if for all possible paths that covers the junction pc has been terminated.
	if consequence.IsTerminated() && alternative.IsTerminated() &&
		solver.Prove(z3.Eq(z3.Or(consequence.pc, alternative.pc), junction.pc)) == nil {
		junction.terminated = true
	}

	// We only merge the variables present in the parent path as the variables declared in the scopes
	// of either the consequence and alternative should not be visible after the corresponding blocks.
	// The variables that have different values in the branches need to be assigned with ITE.
	for _, identifier := range junction.scope.Identifiers() {
		consequenceValue, _ := consequence.scope.Valuation(identifier)
		alternativeValue, _ := alternative.scope.Valuation(identifier)

		if consequenceValue != alternativeValue {
			ite := z3.ITE(consequence.pc, consequenceValue, alternativeValue).Simplify()
			junction.scope.Assign(identifier, ite)
		}
	}

	return junction
}

// Meges an if-then branching.
func (consequence *GoPath) MergeIT() *GoPath {
	path := consequence.MergeITE(consequence.parent)
	return path
}

func (path *GoPath) String() string {
	var builder strings.Builder
	scope := path.scope.String()
	builder.WriteString(scope)
	builder.WriteRune('\n')
	pc := path.pc.String()
	builder.WriteString(fmt.Sprintf("pc=%s", pc))
	return builder.String()
}
