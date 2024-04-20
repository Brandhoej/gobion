package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/automata/language/statements"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_Apply(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)
	x := expressions.NewVariable(symbols.Insert("x"), expressions.IntegerSort)
	y := expressions.NewVariable(symbols.Insert("y"), expressions.IntegerSort)

	update := NewUpdate(
		constraints.NewAssignmentConstraint(
			statements.NewAssignment(x, expressions.NewInteger(1)),
		),
	)

	before := expressions.NewValuationsMap()
	before.Assign(x.Symbol(), expressions.NewInteger(0))
	before.Assign(y.Symbol(), expressions.NewInteger(0))

	valuations := expressions.NewValuationsMap()
	variables := expressions.NewVariablesMap()

	context := z3.NewContext(z3.NewConfig())
	solver := NewConstraintSolver(context.NewSolver(), valuations, variables)

	// Act
	update.Apply(solver)

	// Assert
	t.FailNow()
}
