package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

// a((=2 ∨ =3) ∧ (>=2)) ∧ b(>3 ∨ <=8)
type Guard struct {
	context   *z3.Context
	condition _Expression
}

func NewGuard(context *z3.Context, condition _Expression) Guard {
	return Guard{
		context:   context,
		condition: condition,
	}
}

// Finds the union of all variables and adds conjunctive terms
func (guard Guard) Conjunction(guards ...Guard) Guard {
	conditions := make([]_Expression, 0, len(guards)+1)
	conditions = append(conditions, guard.condition)
	for _, guard := range guards {
		conditions = append(conditions, guard.condition)
	}
	cnf := conditions[0]
	if len(conditions) > 1 {
		cnf = conditions[0].Conjunction(conditions[1:]...)
	}
	return NewGuard(guard.context, cnf)
}

// Finds the intersection of all variables and adds a disjunction term to them.
func (guard Guard) Disjunction(guards ...Guard) Guard {
	condition := make([]_Expression, 0, len(guards)+1)
	condition = append(condition, guard.condition)
	for _, guard := range guards {
		condition = append(condition, guard.condition)
	}
	cnf := condition[0]
	if len(condition) > 1 {
		cnf = condition[0].Disjunction(condition[1:]...)
	}
	return NewGuard(guard.context, cnf)
}

func (guard Guard) Negation() Guard {
	return NewGuard(guard.context, guard.condition.Negation())
}

func (guard Guard) IsSAT(variables Variables) bool {
	solver := newSolver(guard.context.NewSolver())
	solver.Assert(guard.condition)
	return solver.HasSolution()
}

func (guard Guard) IsSatisfiable(variables Variables, valuations Valuations) bool {
	solver := newSolver(guard.context.NewSolver())
	solver.Assert(guard.condition)

	valuations.All(func(symbol symbols.Symbol, value Value) bool {
		if variable, exists := variables.Variable(symbol); exists {
			ast := value.relation.ast(variable, value.ast)
			expression := newExpression(ast)
			solver.Assert(expression)
		} else {
			panic("Unknown variable symbol")
		}

		return true
	})

	return solver.HasSolution()
}

func (guard Guard) String() string {
	return guard.condition.String()
}
