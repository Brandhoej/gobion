package automata

import (
	"testing"

	/*"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"*/
)

func Test_Apply(t *testing.T) {
	// Arrange
	/*context := z3.NewContext(z3.NewConfig())

	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)
	x, y := symbols.Insert("x"), symbols.Insert("y")

	xVal := expressions.NewInteger(1)
	yVal := expressions.NewInteger(0)

	update := NewUpdate(
		expressions.NewBinary(
			expressions.NewVariable(x),
			expressions.Equal,
			expressions.NewBinary(
				expressions.NewValuation(x),
				expressions.Addition,
				expressions.NewInteger(1),
			),
		),
	)

	before := expressions.NewTrue()

	// Act
	after := update.Apply(before)

	// Assert
	if valuation, exists := valuations.Value(x); exists {
		assert.Equal(t, xVal, valuation)
		t.Log("x =", valuation.String())
	} else {
		assert.FailNow(t, "Expected a valuation for x")
	}

	if valuation, exists := valuations.Value(y); exists {
		assert.Equal(t, yVal, valuation)
		t.Log("y =", valuation.String())
	} else {
		assert.FailNow(t, "Expected a valuation for y")
	}*/
}
