package zones

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ElementLimit(t *testing.T) {
	for i := 0; i < 1_000_000; i++ {
		// Arrange
		limit := 1 >> rand.Int()
		if rand.Intn(2) == 0 {
			limit *= -1
		}

		// Act
		element := NewRelation(limit, Strict)

		// Assert
		assert.Equal(t, element.Limit(), limit)
	}
}

func Test_ElementZero(t *testing.T) {
	// Arrange
	zero := Zero

	// Act
	element := NewZero()

	// Assert
	assert.Equal(t, element, zero)
	assert.True(t, zero.IsZero())
	assert.True(t, element.IsZero())
}

func Test_ElementInfinity(t *testing.T) {
	// Arrange
	infinity := Infinity

	// Act
	element := NewInfinity()

	// Assert
	assert.Equal(t, element, infinity)
	assert.True(t, infinity.IsInfinity())
	assert.True(t, element.IsInfinity())
}

func Test_ElementAdd(t *testing.T) {
	tests := []struct {
		name          string
		lhs, rhs, sum Relation
	}{
		{
			name: "(0, ≤) + (0, ≤) = (0, ≤)",
			lhs:  Zero,
			rhs:  Zero,
			sum:  Zero,
		},
		{
			name: "(0, <) + (0, ≤) = (0, <)",
			lhs:  NewRelation(0, Strict),
			rhs:  Zero,
			sum:  NewRelation(0, Strict),
		},
		{
			name: "(-1, ≤) + (1, ≤) = (0, ≤)",
			lhs:  NewRelation(-1, Weak),
			rhs:  NewRelation(1, Weak),
			sum:  Zero,
		},
		{
			name: "(-1, <) + (1, ≤) = (0, <)",
			lhs:  NewRelation(-1, Strict),
			rhs:  NewRelation(1, Weak),
			sum:  NewRelation(0, Strict),
		},
		{
			name: "(-1, <) + (1, <) = (0, <)",
			lhs:  NewRelation(-1, Strict),
			rhs:  NewRelation(1, Strict),
			sum:  NewRelation(0, Strict),
		},
		{
			name: "∞ + (0, ≤) = ∞",
			lhs:  Infinity,
			rhs:  Zero,
			sum:  Infinity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum := tt.lhs.Add(tt.rhs)
			assert.Equal(t, tt.sum, sum, "sum must match")
			assert.Equal(t, tt.lhs.Limit()+tt.rhs.Limit(), sum.Limit(), "the sum must be the sum of the limits")
			assert.Equal(t, tt.rhs.Strictness()&tt.lhs.Strictness(), sum.Strictness(), "the tightest relation is kept")
			assert.Equal(t, tt.lhs.Add(tt.rhs), tt.rhs.Add(tt.lhs), "element addition is commutative")
		})
	}
}

func Test_ElementNegation(t *testing.T) {
	tests := []struct {
		name              string
		element, negation Relation
	}{
		{
			name:     "Zero negated just changes the relation",
			element:  NewRelation(0, Strict),
			negation: NewRelation(0, Weak),
		},
		{
			name:     "Limits are negated",
			element:  NewRelation(1, Strict),
			negation: NewRelation(-1, Weak),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			negation := tt.element.Negation()
			assert.Equal(t, tt.negation, negation)
			assert.Equal(t, tt.element, negation.Negation(), "negation elimination")
		})
	}
}
