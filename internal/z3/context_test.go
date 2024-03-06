package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolUniqueness(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	// Act
	x1 := context.NewConstant(WithName("x"), context.IntegerSort())
	x2 := context.NewConstant(WithName("x"), context.IntegerSort())
	solver.Assert(Not(Eq(x1, x2)))
	sat := solver.Check()

	// Assert
	assert.Equal(t, x1.z3AST, x2.z3AST)
	assert.Equal(t, true, sat.IsFalse())
}

func TestSymbols(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)

	// Act
	constant := context.NewConstant(
		WithInt(0), context.IntegerSort(),
	)

	// Assert
	assert.Equal(t, "k!0", constant.String())
}
