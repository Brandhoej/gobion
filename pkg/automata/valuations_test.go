package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_Value(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[int](
		symbols.NewSymbolsFactory(),
	)
	context := z3.NewContext(z3.NewConfig())
	integer := context.IntegerSort()
	valuations := NewValuationsMap(context)
	a, b, c := symbols.Insert(1), symbols.Insert(2), symbols.Insert(3)
	one, two := context.NewInt(1, integer), context.NewInt(2, integer)

	// Act
	valuations.Assign(a, one, EQ)
	valuations.Assign(b, two, EQ)
	valuations.Assign(a, two, EQ)
	valuations.Assign(c, one, EQ)

	// Assert
	if value, exists := valuations.Value(a); true {
		assert.True(t, exists)
		assert.Equal(t, two, value.ast)
	}
	if value, exists := valuations.Value(b); true {
		assert.True(t, exists)
		assert.Equal(t, two, value.ast)
	}
	if value, exists := valuations.Value(b); true {
		assert.True(t, exists)
		assert.Equal(t, one, value.ast)
	}
}

func Test_SubsetOf(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[int](
		symbols.NewSymbolsFactory(),
	)
	context := z3.NewContext(z3.NewConfig())
	integer := context.IntegerSort()
	valsA := NewValuationsMap(context)
	valsB := NewValuationsMap(context)
	valsC := NewValuationsMap(context)
	a, b, c := symbols.Insert(1), symbols.Insert(2), symbols.Insert(3)
	one, two := context.NewInt(1, integer), context.NewInt(2, integer)

	// {a=1, b<2, c=2}
	valsA.Assign(a, one, EQ)
	valsA.Assign(b, two, LT)
	valsA.Assign(c, two, EQ)
	// {b=1, c=2}
	valsB.Assign(b, one, EQ)
	valsB.Assign(c, two, EQ)
	// {c=2}
	valsC.Assign(c, two, GE)

	// Assert
	assert.True(t, valsA.SubsetOf(valsA))
	assert.False(t, valsA.SubsetOf(valsB))
	assert.False(t, valsA.SubsetOf(valsC))
	assert.True(t, valsB.SubsetOf(valsA))
	assert.True(t, valsB.SubsetOf(valsB))
	assert.False(t, valsB.SubsetOf(valsC))
	assert.True(t, valsC.SubsetOf(valsA))
	assert.True(t, valsC.SubsetOf(valsB))
	assert.True(t, valsC.SubsetOf(valsC))
}