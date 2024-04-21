package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/automata/language/statements"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_Apply(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)
	x := expressions.NewVariable(symbols.Insert("x"))
	y := expressions.NewVariable(symbols.Insert("y"))

	xVal := expressions.NewInteger(1)
	yVal := expressions.NewInteger(0)

	update := NewUpdate(
		constraints.NewAssignmentConstraint(
			statements.NewAssignment(x, xVal),
		),
	)

	variables := expressions.NewVariablesMap()
	valuations := expressions.NewValuationsMap()
	valuations.Assign(x.Symbol(), expressions.NewInteger(0))
	valuations.Assign(y.Symbol(), yVal)

	context := z3.NewContext(z3.NewConfig())
	solver := NewConstraintSolver(context.NewSolver(), variables)

	// Act
	update.Apply(valuations, solver)

	// Assert
	if valuation, exists := valuations.Value(x.Symbol()); exists {
		assert.Equal(t, xVal, valuation)
	} else {
		assert.FailNow(t, "Expected a valuation for x")
	}

	if valuation, exists := valuations.Value(y.Symbol()); exists {
		assert.Equal(t, yVal, valuation)
	} else {
		assert.FailNow(t, "Expected a valuation for y")
	}
}
