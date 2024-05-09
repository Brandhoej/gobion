package symbolic

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
)

type Path struct {
	parent *Path
	scope  Scope
	pc     *z3.AST
}

func NewGlobalPath(scope Scope, pc *z3.AST) *Path {
	return &Path{
		parent: nil,
		scope:  scope,
		pc:     pc,
	}
}

func (path *Path) PC() *z3.AST {
	return path.pc
}

func (path *Path) Solver(function func(solver *z3.Solver)) {
	solver := path.scope.Solver()
	solver.Push()
	solver.Assert(path.pc)
	function(solver)
	solver.Pop(1)
}

func (path *Path) Scope() Scope {
	return path.scope
}

func (path *Path) IsFeasible() (feasible bool) {
	path.Solver(func(solver *z3.Solver) {
		feasible = solver.HasSolution()
	})
	return feasible
}

func (path *Path) Fork(pc *z3.AST) *Path {
	fork := &Path{
		parent: path,
		scope:  path.scope.Branch(),
		pc:     z3.And(path.pc, pc).Simplify(),
	}

	return fork
}

func (path *Path) Join(other *Path) {
	// For each symbol visible in "path".
	// See if there is an updated value for it in "other".
	// If that is the case then we if-then-else the value.
	path.scope.Variables(func(symbol Symbol, _ Scope) {
		valuation := other.scope.Valuation(symbol)
		if valuation == nil {
			return
		}

		path.scope.Assign(
			symbol, z3.ITE(
				// Dont check for nil because the variable is declared in "path".
				other.pc, valuation, path.scope.Valuation(symbol),
			).Simplify(),
		)
	})

	path.pc = z3.Or(path.pc, other.pc).Simplify()
}

func (path *Path) String() string {
	var builder strings.Builder
	path.scope.Variables(func(symbol Symbol, _ Scope) {
		identifier, _ := path.Scope().Identifier(symbol)
		valuation := path.scope.Valuation(symbol)
		builder.WriteString(fmt.Sprintf("%s(k!%v)=%s ", identifier, symbol, valuation.String()))
	})
	builder.WriteString(fmt.Sprintf("\npc=%s", path.pc.String()))
	return builder.String()
}
