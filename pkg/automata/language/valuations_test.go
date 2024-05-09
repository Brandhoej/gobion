package language

import (
	"testing"

	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_Value(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[int](
		symbols.NewSymbolsFactory(),
	)
	a, b, c := symbols.Insert(1), symbols.Insert(2), symbols.Insert(3)
	valuations := NewValuationsMap()
	zero, one := NewInteger(0), NewInteger(1)
	two, three := NewInteger(2), NewInteger(3)

	// Act
	valuations.Assign(a, zero)
	valuations.Assign(b, one)
	valuations.Assign(a, two)
	valuations.Assign(c, three)

	// Assert
	if value, exists := valuations.Value(a); true {
		assert.True(t, exists)
		assert.Equal(t, two, value)
	}
	if value, exists := valuations.Value(b); true {
		assert.True(t, exists)
		assert.Equal(t, one, value)
	}
	if value, exists := valuations.Value(c); true {
		assert.True(t, exists)
		assert.Equal(t, three, value)
	}
}
