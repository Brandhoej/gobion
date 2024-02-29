package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	solver := context.NewSolver()

	// p<->p
	p := context.BoolVar(WithName("p"))
	identity := IFF(p, p)

	solver.Assert(Not(identity))

	// Act
	sat := solver.Check()

	// Assert
	assert.True(t, sat.IsFalse(), "conjecture must be unsatisfiable")
}
