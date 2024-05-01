package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Update struct {
	expression language.Expression
}

func NewUpdate(expression language.Expression) Update {
	return Update{
		expression: expression,
	}
}

func NewEmptyUpdate() Update {
	return Update{
		expression: language.NewTrue(),
	}
}

func (update Update) Apply(expression language.Expression) language.Expression {
	return language.Conjunction(update.expression, expression)
}

func (update Update) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := language.NewPrettyPrinter(&buffer, symbols)
	printer.Expression(update.expression)
	return buffer.String()
}