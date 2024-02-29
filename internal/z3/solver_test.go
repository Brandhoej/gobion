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
	a, b := context.NewConstant(WithName("a"), context.BooleanSort()), context.NewConstant(WithName("b"), context.BooleanSort())

	// Act
	solver.Push()
	implication := Implies(a, b)
	solver.Assert(implication)
	solver.Pop(1)

	// Assert
	assert.Equal(t, "", solver.String())
}
