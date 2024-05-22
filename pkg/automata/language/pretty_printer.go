package language

import (
	"fmt"
	"io"

	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/Brandhoej/gobion/pkg/zones"
)

type PrettyPrinter struct {
	writer  io.Writer
	symbols symbols.Store[any]
}

func NewPrettyPrinter(
	writer io.Writer,
	symbols symbols.Store[any],
) PrettyPrinter {
	return PrettyPrinter{
		writer:  writer,
		symbols: symbols,
	}
}

func (printer PrettyPrinter) WriteString(text string) {
	io.WriteString(printer.writer, text)
}

func (printer PrettyPrinter) Assignment(assignment Assignment) {
	if variable, ok := assignment.lhs.(Variable); ok {
		printer.Variable(variable)
		printer.WriteString("' := ")
		assignment.rhs.Accept(printer)
	}
}

func (printer PrettyPrinter) ClockConstraint(constraint ClockConstraint) {
	lhs, _ := printer.symbols.Item(constraint.lhs)
	rhs, _ := printer.symbols.Item(constraint.rhs)
	printer.WriteString(fmt.Sprintf("%s - %s", lhs, rhs))

	if constraint.relation.Strictness() == zones.Strict {
		printer.WriteString(" < ")
	} else {
		printer.WriteString(" ≤ ")
	}

	if constraint.relation.IsInfinity() {
		printer.WriteString("∞")
	} else {
		printer.WriteString(fmt.Sprintf("%v", constraint.relation.Limit()))
	}
}

func (printer PrettyPrinter) ClockAssignment(assignment ClockAssignment) {
	lhs, _ := printer.symbols.Item(assignment.lhs)
	rhs, _ := printer.symbols.Item(assignment.rhs)
	printer.WriteString(fmt.Sprintf("%s := %s", lhs, rhs))
}

func (printer PrettyPrinter) ClockShift(shift ClockShift) {
	clock, _ := printer.symbols.Item(shift.clock)
	if shift.limit >= 0 {
		printer.WriteString(fmt.Sprintf("%s += %v", clock, shift.limit))
	} else {
		printer.WriteString(fmt.Sprintf("%s -= %v", clock, -shift.limit))
	}
}

func (printer PrettyPrinter) ClockReset(reset ClockReset) {
	clock, _ := printer.symbols.Item(reset.clock)
	printer.WriteString(fmt.Sprintf("%s := %v", clock, reset.limit))
}

func (printer PrettyPrinter) Variable(variable Variable) {
	name, _ := printer.symbols.Item(variable.Symbol())
	printer.WriteString(fmt.Sprintf("%v", name))
}

func (printer PrettyPrinter) Binary(binary Binary) {
	binary.LHS().Accept(printer)
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
	binary.RHS().Accept(printer)
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
	unary.Operand().Accept(printer)
	printer.WriteString(")")
}

func (printer PrettyPrinter) BlockExpression(block BlockExpression) {
	tautology := false
	if boolean, ok := block.expression.(Boolean); ok && boolean.value {
		tautology = true
	}

	if !tautology {
		printer.WriteString("{")
	}

	for idx := range block.statements {
		if idx > 0 {
			printer.WriteString("; ")
		}
		block.statements[idx].Accept(printer)
	}

	if !tautology {
		block.expression.Accept(printer)
		printer.WriteString("}")
	}
}

func (printer PrettyPrinter) IfThenElse(ite IfThenElse) {
	ite.condition.Accept(printer)
	printer.WriteString(" ? ")
	ite.consequence.Accept(printer)
	printer.WriteString(" : ")
	ite.alternative.Accept(printer)
}
