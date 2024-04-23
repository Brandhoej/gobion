package constraints

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_SymbolicInterpretation(t *testing.T) {
	context := z3.NewContext(z3.NewConfig())
	solver := context.NewSolver()

	syms := symbols.NewSymbolsMap[string](symbols.NewSymbolsFactory())
	x := syms.Insert("x")

	xConst := context.NewConstant(
		z3.WithInt(int(x)), context.IntegerSort(),
	)

	zero := context.NewInt(0, context.IntegerSort())
	one := context.NewInt(1, context.IntegerSort())

	variables := expressions.NewVariablesMap[*z3.Sort]()
	variables.Declare(x, context.IntegerSort())
	valuations := expressions.NewValuationsMap[*z3.AST]()
	valuations.Assign(x, zero)

	tests := []struct {
		name       string
		constraint Constraint
		expected   *z3.AST
	}{
		{
			name: "x>1 ∧ x=0",
			constraint: Conjunction(
				NewLogicalConstraint(
					expressions.NewBinary(
						expressions.NewVariable(x),
						expressions.GreaterThan,
						expressions.NewInteger(1),
					),
				),
				NewLogicalConstraint(
					expressions.NewBinary(
						expressions.NewVariable(x),
						expressions.Equal,
						expressions.NewInteger(1),
					),
				),
			),
			expected: z3.And(
				z3.GT(xConst, one),
				z3.Eq(xConst, zero),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valuations := valuations.Copy()
			interpreter := NewSymbolicInterpreter(context, valuations, variables)
			actual := interpreter.Interpret(tt.constraint)
			if !solver.Proven(z3.Eq(actual, tt.expected)) {
				assert.Equal(t, tt.expected.String(), actual.String())
			}
		})
	}
}
