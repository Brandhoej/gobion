package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"
import "runtime"

// Incremental solver, possibly specialized by a particular tactic or logic.
type Solver struct {
	context  *Context
	z3Sovler C.Z3_solver
}

// Create a new solver. This solver is a "combined solver" that internally
// in Z3 uses a non-incremental (solver1) and an incremental solver (solver2)
// This combined solver changes its behaviour based on how it is used
// and how its parameters are set.
//
// If the solver is used in a non incremental way (i.e. no calls to
// Z3_solver_push() or Z3_solver_pop(), and no calls to Z3_solver_assert()
// or Z3_solver_assert_and_track() after checking satisfiability
// without an intervening Z3_solver_reset()) then solver1
// will be used. This solver will apply Z3's "default" tactic.
//
// The "default" tactic will attempt to probe the logic used by the
// assertions and will apply a specialized tactic if one is supported.
// Otherwise the general `(and-then simplify smt)` tactic will be used.
//
// If the solver is used in an incremental way then the combined solver
// will switch to using solver2 (which behaves similarly to the general
// "smt" tactic).
//
// Note however it is possible to set the solver2_timeout,
// solver2_unknown, and ignore_solver1 parameters of the combined
// solver to change its behaviour.
func (context *Context) NewSolver() (solver *Solver) {
	context.do(func() {
		solver = &Solver{
			context:  context,
			z3Sovler: C.Z3_mk_solver(context.z3Context),
		}

		// User must use Z3_solver_inc_ref and Z3_solver_dec_ref to manage solver objects.
		// Even if the context was created using Z3_mk_context instead of Z3_mk_context_rc.
		C.Z3_solver_inc_ref(context.z3Context, solver.z3Sovler)
	})

	runtime.SetFinalizer(solver, func(solver *Solver) {
		context.do(func() {
			C.Z3_solver_dec_ref(context.z3Context, solver.z3Sovler)
		})
	})

	return solver
}

func (solver *Solver) Context() *Context {
	return solver.context
}

func (solver *Solver) Assert(ast *AST) {
	solver.context.do(func() {
		C.Z3_solver_assert(solver.context.z3Context, solver.z3Sovler, ast.z3AST)
	}, solver, ast)
}

func (solver *Solver) ReasonUnknown() (reason string) {
	return compute(solver.context, func() string {
		return C.GoString(
			C.Z3_solver_get_reason_unknown(solver.context.z3Context, solver.z3Sovler),
		)
	}, solver)
}

func (solver *Solver) Check() LiftedBoolean {
	return compute(solver.context, func() LiftedBoolean {
		return LiftedBoolean(
			C.Z3_solver_check(solver.context.z3Context, solver.z3Sovler),
		)
	}, solver)
}

func (solver *Solver) Reset() {
	solver.context.do(func() {
		C.Z3_solver_reset(solver.context.z3Context, solver.z3Sovler)
	}, solver)
}

func (solver *Solver) Push() {
	solver.context.do(func() {
		C.Z3_solver_push(solver.context.z3Context, solver.z3Sovler)
	}, solver)
}

func (solver *Solver) Pop(amount uint32) {
	solver.context.do(func() {
		amount := C.uint(amount)
		C.Z3_solver_pop(solver.context.z3Context, solver.z3Sovler, amount)
	}, solver)
}

func (solver *Solver) Depth() uint {
	return compute(solver.context, func() uint {
		return uint(C.Z3_solver_get_num_scopes(solver.context.z3Context, solver.z3Sovler))
	}, solver)
}

func (solver *Solver) String() string {
	return compute(solver.context, func() string {
		return C.GoString(
			C.Z3_solver_to_string(solver.context.z3Context, solver.z3Sovler),
		)
	}, solver)
}

func (solver *Solver) Dimacs(includeNames bool) string {
	return compute(solver.context, func() string {
		return C.GoString(
			C.Z3_solver_to_dimacs_string(solver.context.z3Context, solver.z3Sovler, C.bool(includeNames)),
		)
	}, solver)
}

func (solver *Solver) Prove(proposition *AST) *Model {
	solver.Push()
	defer solver.Pop(1)

	contradiction := Not(proposition)
	solver.Assert(contradiction)
	sat := solver.Check()
	if sat.IsTrue() {
		return solver.Model()
	}

	return nil
}

func (solver *Solver) Proven(proposition *AST) bool {
	solver.Push()
	defer solver.Pop(1)

	contradiction := Not(proposition)
	solver.Assert(contradiction)
	return solver.Check().IsFalse()
}

func (solver *Solver) IsTautology(proposition *AST) *Model {
	return solver.Prove(
		Eq(proposition, solver.True()),
	)
}

func (solver *Solver) IsContradiction(proposition *AST) *Model {
	return solver.Prove(
		Eq(proposition, solver.False()),
	)
}

func (solver *Solver) HasSolutionFor(proposition *AST) bool {
	solver.Push()
	defer solver.Pop(1)
	solver.Assert(proposition)
	return solver.Check().IsTrue()
}

func (solver *Solver) HasSolution() bool {
	return solver.Check().IsTrue()
}

func (solver *Solver) True() *AST {
	return solver.context.NewTrue()
}

func (solver *Solver) False() *AST {
	return solver.context.NewFalse()
}
