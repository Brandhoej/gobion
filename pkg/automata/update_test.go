package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_Apply(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)
	variables := language.NewVariablesMap()
	x := variables.Declare(symbols.Insert("x"), language.IntegerSort)
	y := variables.Declare(symbols.Insert("y"), language.IntegerSort)

	update := NewUpdate(
		language.NewAssignment(x, language.NewInteger(1)),
	)

	before := language.NewValuationsMap()
	before.Assign(x.Symbol(), language.NewInteger(0))
	before.Assign(y.Symbol(), language.NewInteger(0))
	
	context := z3.NewContext(z3.NewConfig())
	statements := language.NewSymbolicStatementInterpreter(context, before)

	// Act
	after := update.Apply(statements, before)

	// Assert
	t.Log(after)
	t.FailNow()
}
