package constraints

import (
	"io"

	"github.com/Brandhoej/gobion/pkg/automata/language/statements"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type PrettyPrinter struct {
	writer io.Writer
	statements statements.PrettyPrinter
}

func NewPrettyPrinter(
	writer io.Writer,
	symbols symbols.Store[any],
) PrettyPrinter {
	return PrettyPrinter{
		writer: writer,
		statements: statements.NewPrettyPrinter(writer, symbols),
	}
}

func (printer PrettyPrinter) WriteString(text string) {
	io.WriteString(printer.writer, text)
}

func (printer PrettyPrinter) Constraint(constraint Constraint) {
	switch cast := any(constraint).(type) {
	case ExpressionConstraint:
		printer.ExpressionConstraint(cast)
	case BinaryConstraint:
		printer.BinaryConstraint(cast)
	case UnaryConstraint:
		printer.UnaryConstraint(cast)
	}
}

func (printer PrettyPrinter) ExpressionConstraint(constraint ExpressionConstraint) {
	printer.statements.Expression(constraint.expression)
}

func (printer PrettyPrinter) BinaryConstraint(constraint BinaryConstraint) {
	printer.Constraint(constraint.lhs)
	switch constraint.operator {
	case LogicalAnd:
		printer.WriteString("&&")
	case LogicalOr:
		printer.WriteString("||")
	}
	printer.Constraint(constraint.rhs)
}

func (printer PrettyPrinter) UnaryConstraint(constraint UnaryConstraint) {
	switch constraint.operator {
	case LogicalNegation:
		printer.WriteString("!")
	}
	printer.WriteString("(")
	printer.Constraint(constraint.operand)
	printer.WriteString(")")
}