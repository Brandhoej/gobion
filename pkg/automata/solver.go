package automata

import "github.com/Brandhoej/gobion/internal/z3"

type Solver struct {
	backing *z3.Solver
}

func newSolver(solver *z3.Solver) Solver {
	return Solver{
		backing: solver,
	}
}

func (solver Solver) Assert(expression _Expression) {
	solver.backing.Assert(expression.ast)
}

func (solver Solver) HasSolution() bool {
	return solver.backing.HasSolution()
}
