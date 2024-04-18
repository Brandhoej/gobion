package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()
	a := context.NewConstant(WithName("a"), context.BooleanSort())
	b := context.NewConstant(WithName("b"), context.BooleanSort())

	// Act
	solver.Push()
	implication := Implies(a, b)
	solver.Assert(implication)
	solver.Pop(1)

	// Assert
	assert.Equal(t, "", solver.String())
}

func TestMultipleAssertions(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()
	a := context.NewConstant(WithName("a"), context.IntegerSort())
	one := context.NewInt(1, context.IntegerSort())
	two := context.NewInt(2, context.IntegerSort())
	three := context.NewInt(3, context.IntegerSort())

	// Act
	solver.Assert(Eq(a, one))
	solver.Assert(Eq(a, two))
	solver.Assert(Eq(a, three))

	// Assert
	assert.False(t, solver.HasSolution())
}

func TestAsdDsa(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()
	a := context.NewConstant(WithName("a"), context.IntegerSort())
	b := context.NewConstant(WithName("b"), context.IntegerSort())
	p := context.NewConstant(WithName("p"), context.BooleanSort())
	one := context.NewInt(1, context.IntegerSort())
	two := context.NewInt(2, context.IntegerSort())
	// three := context.NewInt(3, context.IntegerSort())

	// Act
	// solver.Assert(context.NewFalse())
	solver.Assert(Eq(one, ITE(p, one, two)))
	solver.Assert(Eq(a, b))

	// Assert
	assert.False(t, solver.HasSolution())
}
