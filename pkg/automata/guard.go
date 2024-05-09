package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Guard struct {
	condition language.Expression
}

func NewGuard(condition language.Expression) Guard {
	return Guard{
		condition: condition,
	}
}

func NewTrueGuard() Guard {
	return NewGuard(language.NewTrue())
}

func NewFalseGuard() Guard {
	return NewGuard(language.NewFalse())
}

// Finds the union of all variables and adds conjunctive terms
func (guard Guard) Conjunction(guards ...Guard) Guard {
	conditions := make([]language.Expression, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].condition
	}
	conjunction := language.Conjunction(guard.condition, conditions...)
	return NewGuard(conjunction)
}

// Finds the intersection of all variables and adds a disjunction term to them.
func (guard Guard) Disjunction(guards ...Guard) Guard {
	conditions := make([]language.Expression, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].condition
	}
	conjunction := language.Disjunction(guard.condition, conditions...)
	return NewGuard(conjunction)
}

func (guard Guard) Negation() Guard {
	negation := language.LogicalNegate(guard.condition)
	return NewGuard(negation)
}

func (guard Guard) IsSatisfied(valuations language.Valuations, solver *Interpreter) bool {
	return solver.IsSatisfied(valuations, guard.condition)
}

func (guard Guard) IsSatisfiable(solver *Interpreter) bool {
	return solver.IsSatisfiable(guard.condition)
}

func (guard Guard) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := language.NewPrettyPrinter(&buffer, symbols)
	guard.condition.Accept(printer)
	return buffer.String()
}
