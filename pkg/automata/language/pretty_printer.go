package language

import (
	"fmt"
	"io"

	"github.com/Brandhoej/gobion/pkg/symbols"
)

type PrettyPrinter struct {
	writer io.Writer
	symbols symbols.Store[any]
}

func NewPrettyPrinter(
	writer io.Writer,
	symbols symbols.Store[any],
) PrettyPrinter {
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
	case Assignment:
		printer.Assignment(cast)
	default:
		panic("Unknown statement type")
	}
}

func (printer PrettyPrinter) Assignment(assignment Assignment) {
	printer.Variable(assignment.variable)
	printer.WriteString("' := ")
	printer.Expression(assignment.valuation)
}

func (printer PrettyPrinter) Expression(expression Expression) {
	switch cast := any(expression).(type) {
	case Variable:
		printer.Variable(cast)
	case Binary:
		printer.Binary(cast)
	case Integer:
		printer.Integer(cast)
	case Boolean:
		printer.Boolean(cast)
	case Unary:
		printer.Unary(cast)
	case BlockExpression:
		printer.BlockExpression(cast)
	default:
		panic("Unknown expression type")
	}
}

func (printer PrettyPrinter) Variable(variable Variable) {
	name, _ := printer.symbols.Item(variable.Symbol())
	printer.WriteString(fmt.Sprintf("%v", name))
}

func (printer PrettyPrinter) Binary(binary Binary) {
	printer.Expression(binary.LHS())
	switch binary.Operator() {
	case Equal:
		printer.WriteString(" = ")
	case NotEqual:
		printer.WriteString(" ≠ ")
	case LessThan:
		printer.WriteString(" < ")
	case LessThanEqual:
		printer.WriteString(" ≤ ")
	case GreaterThan:
		printer.WriteString(" = ")
	case GreaterThanEqual:
		printer.WriteString(" ≥ ")
	case LogicalAnd:
		printer.WriteString(" ∧ ")
	case LogicalOr:
		printer.WriteString(" ∨ ")
	case Addition:
		printer.WriteString(" + ")
	case Subtraction:
		printer.WriteString(" - ")
	case Implication:
		printer.WriteString(" → ")
	default:
		panic("Unknown binary operator")
	}
	printer.Expression(binary.RHS())
}

func (printer PrettyPrinter) Integer(integer Integer) {
	printer.WriteString(
		fmt.Sprintf("%v", integer.Value()),
	)
}

func (printer PrettyPrinter) Boolean(boolean Boolean) {
	printer.WriteString(
		fmt.Sprintf("%v", boolean.Value()),
	)
}

func (printer PrettyPrinter) Unary(unary Unary) {
	switch unary.Operator() {
	case LogicalNegation:
		printer.WriteString("¬")
	default:
		panic("Unknown unary operator")
	}
	printer.WriteString("(")
	printer.Expression(unary.Operand())
	printer.WriteString(")")
}

func (printer PrettyPrinter) BlockExpression(block BlockExpression) {
	printer.WriteString("{")
	for idx := range block.statements {
		printer.Statement(block.statements[idx])
		printer.WriteString("; ")
	}
	printer.Expression(block.expression)
	printer.WriteString("}")
}