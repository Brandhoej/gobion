package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Guard struct {
	constraint constraints.Constraint
}

func NewGuard(constraint constraints.Constraint) Guard {
	return Guard{
		constraint: constraint,
	}
}

func NewTrueGuard() Guard {
	return NewGuard(constraints.NewTrue())
}

func NewFalseGuard() Guard {
	return NewGuard(constraints.NewFalse())
}

// Finds the union of all variables and adds conjunctive terms
func (guard Guard) Conjunction(guards ...Guard) Guard {
	conditions := make([]constraints.Constraint, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].constraint
	}
	conjunction := constraints.Conjunction(guard.constraint, conditions...)
	return NewGuard(conjunction)
}

// Finds the intersection of all variables and adds a disjunction term to them.
func (guard Guard) Disjunction(guards ...Guard) Guard {
	conditions := make([]constraints.Constraint, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].constraint
	}
	disjunction := constraints.Disjunction(guard.constraint, conditions...)
	return NewGuard(disjunction)
}

func (guard Guard) Negation() Guard {
	negation := constraints.LogicalNegate(guard.constraint)
	return NewGuard(negation)
}

func (guard Guard) IsSatisfiable(valuations expressions.Valuations[*z3.AST], solver *ConstraintSolver) bool {
	return solver.Satisfies(valuations, guard.constraint)
}

func (guard Guard) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := constraints.NewPrettyPrinter(&buffer, symbols)
	printer.Constraint(guard.constraint)
	return buffer.String()
}
