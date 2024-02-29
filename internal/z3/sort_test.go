package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameAs(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)
	tests := []struct {
		lhs  *Sort
		rhs  *Sort
		same bool
	}{
		{
			lhs:  context.IntegerSort(),
			rhs:  context.IntegerSort(),
			same: true,
		},
		{
			lhs:  context.BooleanSort(),
			rhs:  context.BooleanSort(),
			same: true,
		},
		{
			lhs:  context.BooleanSort(),
			rhs:  context.IntegerSort(),
			same: false,
		},
	}

	for _, test := range tests {
		// Act
		same := test.lhs.SameAs(test.rhs)

		// Assert
		assert.Equal(t, test.same, same)
	}
}
