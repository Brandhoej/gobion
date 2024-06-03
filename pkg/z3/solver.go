package z3

import "github.com/Brandhoej/gobion/internal/z3"

type Solver struct {
	_solver *z3.Solver
	backlock []*AST
}

func newSolver(_solver *z3.Solver) *Solver {
	return &Solver{
		_solver: _solver,
		backlock: make([]*AST, 0, 16),
	}
}

func (solver *Solver) Check() LiftedBoolean {
	return LiftedBoolean(solver._solver.Check())
}
