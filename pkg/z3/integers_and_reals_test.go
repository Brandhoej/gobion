package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdditionCommutative(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	a, b := context.IntVar(WithName("a")), context.IntVar(WithName("b"))

	cummutative := Eq(Add(a, b), Add(b, a))
	solver.Assert(Not(cummutative))

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}

func TestAdditionAssociative(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	a, b, c := context.IntVar(WithName("a")), context.IntVar(WithName("b")), context.IntVar(WithName("c"))
	cummutative := Eq(Add(Add(a, b), c), Add(a, Add(b, c)))
	solver.Assert(Not(cummutative))

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}

func TestAdditionModel(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)

	a, b, c := context.IntVar(WithName("a")), context.IntVar(WithName("b")), context.IntVar(WithName("c"))

	// Act
	cummutative := Eq(Add(Add(a, b), c), Add(a, Add(b, c)))

	// Assert
	assert.Equal(t, "(= (+ a b c) (+ a b c))", cummutative.String())
}
