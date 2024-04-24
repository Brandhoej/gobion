package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Update struct {
	expression expressions.Expression
}

func NewUpdate(expression expressions.Expression) Update {
	return Update{
		expression: expression,
	}
}

func NewEmptyUpdate() Update {
	return Update{
		expression: expressions.NewTrue(),
	}
}

func (update Update) Apply(expression expressions.Expression) expressions.Expression {
	return expressions.Conjunction(update.expression, expression)
}

func (update Update) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := expressions.NewPrettyPrinter(&buffer, symbols)
	printer.Expression(update.expression)
	return buffer.String()
}