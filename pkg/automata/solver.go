package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
)

type ConstraintSolver struct {
	interpreter constraints.SymbolicInterpreter
	constraints []constraints.Constraint // TODO: Make a rooted tree.
	backing     *z3.Solver
}

func NewConstraintSolver(solver *z3.Solver, valuations expressions.Valuations) *ConstraintSolver {
	return &ConstraintSolver{
		interpreter: constraints.NewSymbolicInterpreter(
			solver.Context(), valuations,
		),
		backing: solver,
	}
}

func (solver *ConstraintSolver) Branch() *ConstraintSolver {
	constraints := make([]constraints.Constraint, len(solver.constraints))
	copy(constraints, solver.constraints)
	return &ConstraintSolver{
		interpreter: solver.interpreter,
		backing:     solver.backing,
		constraints: constraints,
	}
}

func (solver *ConstraintSolver) setup() {
	for idx := range solver.constraints {
		proposition := solver.interpreter.Interpret(
			solver.constraints[idx],
		)
		solver.backing.Assert(proposition)
	}
}

func (solver *ConstraintSolver) Assert(constraint constraints.Constraint) {
	solver.constraints = append(solver.constraints, constraint)
}

func (solver *ConstraintSolver) HasSolutionFor(constraint constraints.Constraint) bool {
	solver.backing.Push()
	defer solver.backing.Pop(1)
	solver.setup()

	proposition := solver.interpreter.Interpret(constraint)
	return solver.backing.HasSolutionFor(proposition)
}

func (solver *ConstraintSolver) HasSolution() bool {
	solver.backing.Push()
	defer solver.backing.Pop(1)
	solver.setup()

	return solver.backing.HasSolution()
}
