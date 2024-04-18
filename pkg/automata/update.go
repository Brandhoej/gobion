package automata

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Update struct {
	context *z3.Context
	/* A sequence of assignments where order is kept.
	 * For each assignment we have the symbol and an expression.
	 * When an assignment is done we substitute all variables in the expression
	 * with values from a valuation set when this has been done the result is
	 * simplified and assigned in the valuation set to the symbol. */
	assignments []Assignment
}

func NewUpdate(context *z3.Context, assignments ...Assignment) Update {
	return Update{
		context: context,
		assignments: assignments,
	}
}

func (update Update) Apply(variables Variables, valuations Valuations, forward bool) Valuations {
	result := NewValuationsMap(update.context)

	indices := make(map[symbols.Symbol]int)
	vals := make([]*z3.AST, 0)
	
	// Step 1: Copy and store all valuations.
	valuations.All(func(symbol symbols.Symbol, value Value) bool {
		result.Assign(symbol, value.ast, value.relation)

		index := len(vals)
		vals = append(vals, value.ast)
		indices[symbol] = index

		return true
	})

	// Step 2: Store all variables that we have a valuation for.
	vars := make([]*z3.AST, len(vals))
	variables.All(func(symbol symbols.Symbol, variable *z3.AST) bool {
		if index, hasValue := indices[symbol]; hasValue {
			vars[index] = variable
		}
		return true
	})

	// Step 3: Apply all assignments with the use of substitution.
	for idx := range update.assignments {
		if !forward {
			idx = len(update.assignments)-(idx + 1)
		}
		assignment := update.assignments[idx]
		symbol := assignment.symbol
		expression := assignment.expression
		substitution := expression.Substitute(vars, vals)
		result.Assign(symbol, substitution, assignment.operator)
	}

	return result
}

func (update Update) String() string {
	var builder strings.Builder
	for _, assignment := range update.assignments {
		builder.WriteString(
			fmt.Sprintf("k!%v %v %v", assignment.symbol, assignment.operator.String(), assignment.expression.String()),
		)
	}

	return builder.String()
}