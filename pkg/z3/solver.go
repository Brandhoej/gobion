package z3

import "github.com/Brandhoej/gobion/internal/z3"

type Solver struct {
	_solver *z3.Solver
}

func (solver *Solver) Assert(valuation valuation) {
	solver._solver.Assert(valuation.ast())
}

func (solver *Solver) Check() LiftedBoolean {
	return LiftedBoolean(solver._solver.Check())
}
