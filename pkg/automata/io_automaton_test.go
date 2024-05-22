package automata

/*import (
	"bytes"
	"testing"

	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/Brandhoej/gobion/pkg/zones"
)

func Test_Administration(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[any](symbols.NewSymbolsFactory())
	grant := Action(symbols.Insert("grant"))
	publication := Action(symbols.Insert("publication"))
	news := Action(symbols.Insert("news"))
	coin := Action(symbols.Insert("coin"))
	reference := symbols.Insert("0")
	z := symbols.Insert("z")

	builder := NewIOAutomatonBuilder()
	builder.AddInputs(grant, publication)
	builder.AddOutputs(coin, news)
	tl := builder.AddInitial("tl")
	tr := builder.AddLocation("tr", WithInvariant(
		NewInvariant(
			language.NewClockCondition(z, reference, zones.NewRelation(2, zones.Strict)),
		),
	))
	bl := builder.AddLocation("bl", WithInvariant(
		NewInvariant(
			language.NewClockCondition(z, reference, zones.NewRelation(2, zones.Strict)),
		),
	))
	br := builder.AddLocation("br")
	builder.AddEdge(tl, grant, tr, WithUpdate(
		NewUpdate(
			language.NewBlockExpression(
				language.NewTrue(), language.NewClockReset(z, 0),
			),
		),
	))
	builder.AddEdge(tl, publication, bl, WithUpdate(
		NewUpdate(
			language.NewBlockExpression(
				language.NewTrue(), language.NewClockReset(z, 0),
			),
		),
	))
	builder.AddLoop(tr, grant)
	builder.AddLoop(tr, publication)
	builder.AddEdge(tr, coin, br)
	builder.AddEdge(bl, news, tl)
	builder.AddLoop(bl, grant)
	builder.AddLoop(bl, publication)
	builder.AddEdge(br, publication, bl, WithUpdate(
		NewUpdate(
			language.NewBlockExpression(
				language.NewTrue(), language.NewClockReset(z, 0),
			),
		),
	))
	builder.AddLoop(br, grant)
	automaton := builder.Build()

	// Act
	var buffer bytes.Buffer
	automaton.DOT(&buffer, symbols)

	// Assert
	t.Log(buffer.String())
	t.FailNow()
}
*/