package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTautology(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	p := context.BoolVar(WithName("p"))
	solver.Assert(Not(Or(p, Not(p))))

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}

func TestTransitivity(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	// If p->q and q->r then p->r
	p, q, r := context.BoolVar(WithName("p")), context.BoolVar(WithName("q")), context.BoolVar(WithName("r"))
	transitivity := Implies(And(Implies(p, q), Implies(q, r)), Implies(p, r))

	solver.Assert(Not(transitivity))

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}

func TestDemorganProof(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	// !(x && y) == (!x || !y)
	x, y := context.BoolVar(WithName("x")), context.BoolVar(WithName("y"))
	l := Not(And(x, y))
	r := Or(Not(x), Not(y))
	demorgan := IFF(l, r)
	conjecture := Not(demorgan)

	// !(!(x && y) == (!x || !y))
	solver.Assert(conjecture)

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}

// Testing equivalence (IFF) no specific values are assigned to the variables x and y.
// The solver does not make any concrete assignments resulting in an empty model string.
func TestDemorganModel(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	// !(x && y) == (!x || !y)
	x, y := context.BoolVar(WithName("x")), context.BoolVar(WithName("y"))
	l := Not(And(x, y))
	r := Or(Not(x), Not(y))
	demorgan := IFF(l, r)

	solver.Assert(demorgan)
	isSAT := solver.Check()

	// Act
	model := solver.Model()
	text := model.String()

	// Assert
	assert.True(t, isSAT.IsTrue(), "demorgan must have a solution")
	assert.Equal(t, "", text, "No conrete valuations when using equivalence")
}

func TestModulPonens(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	premise, conclusion := context.BoolVar(WithName("p")), context.BoolVar(WithName("q"))
	modusPonen := Implies(And(Implies(premise, conclusion), premise), conclusion)

	// Modus ponens: If p implies q then if p is true q is true.
	// Can we find an example where !p->q and q is true?
	solver.Assert(Not(modusPonen))

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}

func TestDoubleNegationElimination(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	// If p->q and q->r then p->r
	p := context.BoolVar(WithName("p"))
	doubleNegation := IFF(Not(Not(p)), p)

	solver.Assert(Not(doubleNegation))

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}
