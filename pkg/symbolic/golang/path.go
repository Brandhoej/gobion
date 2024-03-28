package golang

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
)

type GoPath struct {
	parent  *GoPath
	scope   *GoScope
	context *z3.Context
	pc      *z3.AST
}

func NewGoGlobalPath(context *z3.Context) *GoPath {
	return &GoPath{
		parent:  nil,
		scope:   NewGoGlobalScope(0, context.NewSolver()),
		context: context,
		pc:      context.NewTrue(),
	}
}

func (path *GoPath) Solver(function func(solver *z3.Solver)) {
	solver := path.scope.State()
	solver.Push()
	solver.Assert(path.pc)
	function(solver)
	solver.Pop(1)
}

func (path *GoPath) IsTautologhy() bool {
	isTautologhy := false
	path.Solver(func(solver *z3.Solver) {
		equality := z3.Eq(path.pc, path.context.NewTrue())
		isTautologhy = solver.Prove(equality) != nil
	})
	return isTautologhy
}

func (path *GoPath) IsFeasible() bool {
	isFeasible := false
	path.Solver(func(solver *z3.Solver) {
		isFeasible = solver.HasSolution()
	})
	return isFeasible
}

func (path *GoPath) IsInfeasible() bool {
	return !path.IsFeasible()
}

func (path *GoPath) Enclose(scopeID int) *GoPath {
	return path.Branch(path.pc, scopeID)
}

func (path *GoPath) EncloseInto(scope *GoScope) *GoPath {
	return path.BranchInto(path.pc, scope)
}

func (path *GoPath) Branch(condition *z3.AST, scopeID int) *GoPath {
	return path.BranchInto(condition, path.scope.Child(scopeID))
}

func (path *GoPath) BranchInto(condition *z3.AST, scope *GoScope) *GoPath {
	branch := &GoPath{
		parent:  path,
		scope:   scope,
		context: path.context,
		pc:      z3.And(path.pc, condition).Simplify(),
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
		parent:  parent,
		scope:   parent.scope,
		context: parent.context,
		pc:      parent.pc,
	}

	// We only merge the variables present in the parent path as the variables declared in the scopes
	// of either the consequence and alternative should not be visible after the corresponding blocks.
	// The variables that have different values in the branches need to be assigned with ITE.
	junction.scope.Traverse(func(scope *GoScope) {
		for _, symbol := range scope.symbols.Symbols() {
			consequenceValue, hasConsequence := consequence.scope.ValuationFor(symbol)
			alternativeValue, hasAlternative := alternative.scope.ValuationFor(symbol)

			if hasConsequence && hasAlternative && consequenceValue != alternativeValue {
				ite := z3.ITE(consequence.pc, consequenceValue, alternativeValue).Simplify()
				junction.scope.AssignTo(symbol, ite)
			}
		}
	})

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
