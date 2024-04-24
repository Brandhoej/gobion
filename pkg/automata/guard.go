package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Guard struct {
	condition expressions.Expression
}

func NewGuard(condition expressions.Expression) Guard {
	return Guard{
		condition: condition,
	}
}

func NewTrueGuard() Guard {
	return NewGuard(expressions.NewTrue())
}

func NewFalseGuard() Guard {
	return NewGuard(expressions.NewFalse())
}

// Finds the union of all variables and adds conjunctive terms
func (guard Guard) Conjunction(guards ...Guard) Guard {
	conditions := make([]expressions.Expression, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].condition
	}
	conjunction := expressions.Conjunction(guard.condition, conditions...)
	return NewGuard(conjunction)
}

// Finds the intersection of all variables and adds a disjunction term to them.
func (guard Guard) Disjunction(guards ...Guard) Guard {
	conditions := make([]expressions.Expression, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].condition
	}
	conjunction := expressions.Disjunction(guard.condition, conditions...)
	return NewGuard(conjunction)
}

func (guard Guard) Negation() Guard {
	negation := expressions.LogicalNegate(guard.condition)
	return NewGuard(negation)
}

func (guard Guard) IsSatisfied(valuations expressions.Valuations[*z3.AST], solver *Interpreter) bool {
	return solver.IsSatisfied(valuations, guard.condition)
}

func (guard Guard) IsSatisfiable(solver *Interpreter) bool {
	return solver.IsSatisfiable(guard.condition)
}

func (guard Guard) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := expressions.NewPrettyPrinter(&buffer, symbols)
	printer.Expression(guard.condition)
	return buffer.String()
}
