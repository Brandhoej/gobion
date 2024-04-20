package statements

import (
	"fmt"
	"io"

	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type PrettyPrinter struct {
	writer io.Writer
	symbols symbols.Store[any]
}

func NewPrettyPrinter(writer io.Writer, symbols symbols.Store[any]) PrettyPrinter {
	return PrettyPrinter{
		writer: writer,
		symbols: symbols,
	}
}

func (printer PrettyPrinter) WriteString(text string) {
	io.WriteString(printer.writer, text)
}

func (printer PrettyPrinter) Statement(statement Statement) {
	switch cast := any(statement).(type) {
	case Block:
		printer.Block(cast)
	case Assignment:
		printer.Assignment(cast)
	default:
		panic("Unknown statement type")
	}
}

func (printer PrettyPrinter) Expression(expression expressions.Expression) {
	switch cast := any(expression).(type) {
	case expressions.Variable:
		printer.Variable(cast)
	case expressions.Binary:
		printer.Binary(cast)
	case expressions.Integer:
		printer.Integer(cast)
	case expressions.Boolean:
		printer.Boolean(cast)
	case expressions.Unary:
		printer.Unary(cast)
	default:
		panic("Unknown expression type")
	}
}

func (printer PrettyPrinter) Variable(variable expressions.Variable) {
	name, _ := printer.symbols.Item(variable.Symbol())
	printer.WriteString(fmt.Sprintf("%v", name))
}

func (printer PrettyPrinter) Binary(binary expressions.Binary) {
	printer.Expression(binary.LHS())
	switch binary.Operator() {
	case expressions.Equal:
		printer.WriteString(" = ")
	case expressions.NotEqual:
		printer.WriteString(" ≠ ")
	case expressions.LessThan:
		printer.WriteString(" < ")
	case expressions.LessThanEqual:
		printer.WriteString(" ≤ ")
	case expressions.GreaterThan:
		printer.WriteString(" = ")
	case expressions.GreaterThanEqual:
		printer.WriteString(" ≥ ")
	case expressions.LogicalAnd:
		printer.WriteString(" ∧ ")
	case expressions.LogicalOr:
		printer.WriteString(" ∨ ")
	case expressions.Addition:
		printer.WriteString(" + ")
	case expressions.Subtraction:
		printer.WriteString(" - ")
	default:
		panic("Unknown binary operator")
	}
	printer.Expression(binary.RHS())
}

func (printer PrettyPrinter) Integer(integer expressions.Integer) {
	printer.WriteString(
		fmt.Sprintf("%v", integer.Value()),
	)
}

func (printer PrettyPrinter) Boolean(boolean expressions.Boolean) {
	printer.WriteString(
		fmt.Sprintf("%v", boolean.Value()),
	)
}

func (printer PrettyPrinter) Unary(unary expressions.Unary) {
	switch unary.Operator() {
	case expressions.LogicalNegation:
		printer.WriteString("¬")
	default:
		panic("Unknown unary operator")
	}
	printer.Expression(unary.Operand())
}

func (printer PrettyPrinter) Block(block Block) {
	printer.WriteString("{")
	for idx, statement := range block.Statements() {
		if idx > 0 {
			printer.WriteString(" ")
		}
		printer.Statement(statement)
	}
	printer.WriteString("}")
}

func (printer PrettyPrinter) Assignment(assignment Assignment) {
	printer.Expression(assignment.LHS())
	printer.WriteString("' = ")
	printer.Expression(assignment.RHS())
}
