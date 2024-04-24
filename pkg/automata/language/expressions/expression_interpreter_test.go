package expressions

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/state"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_SymbolicInterpretation(t *testing.T) {
	context := z3.NewContext(z3.NewConfig())
	solver := context.NewSolver()

	symbols := symbols.NewSymbolsMap[string](symbols.NewSymbolsFactory())
	x, y := symbols.Insert("x"), symbols.Insert("y")
	xConst := context.NewConstant(
		z3.WithInt(int(x)), context.IntegerSort(),
	)

	variables := state.NewVariablesMap[*z3.Sort]()
	variables.Declare(x, context.IntegerSort())
	variables.Declare(y, context.IntegerSort())

	tests := []struct {
		name       string
		expression Expression
		expected   *z3.AST
		valuations func(t *testing.T, valuations state.Valuations[*z3.AST])
	}{
		{
			name:       "true",
			expression: NewBoolean(true),
			expected:   context.NewTrue(),
		},
		{
			name:       "false",
			expression: NewBoolean(false),
			expected:   context.NewFalse(),
		},
		{
			name:       "x",
			expression: NewVariable(x),
			expected:   xConst,
		},
		{
			name:       "y",
			expression: NewVariable(y),
			expected:   context.NewInt(1, context.IntegerSort()),
		},
		{
			name:       "x+y",
			expression: NewBinary(NewVariable(x), Addition, NewVariable(y)),
			expected: z3.Add(
				xConst, context.NewInt(1, context.IntegerSort()),
			),
		},
		{
			name: "x>0 ? 1 : 2",
			expression: NewIfThenElse(
				NewBinary(NewVariable(x), GreaterThan, NewInteger(0)),
				NewInteger(1),
				NewInteger(2),
			),
			expected: context.NewInt(2, context.IntegerSort()),
		},
		{
			name:       "x=2",
			expression: NewBinary(NewVariable(x), Equal, NewInteger(2)),
			expected:   context.NewTrue(),
			valuations: func(t *testing.T, valuations state.Valuations[*z3.AST]) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
		{
			name: "y=3 ∧ x'=2",
			expression: Conjunction(
				NewBinary(NewVariable(y), Equal, NewInteger(3)),
				NewAssignment(NewVariable(x), NewInteger(2)),
			),
			expected: context.NewFalse(),
			valuations: func(t *testing.T, valuations state.Valuations[*z3.AST]) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
		{
			name: "y=1 ∨ x'=2",
			expression: Disjunction(
				NewBinary(NewVariable(y), Equal, NewInteger(1)),
				NewAssignment(NewVariable(x), NewInteger(2)),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations state.Valuations[*z3.AST]) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
		{
			name: "y=1 ∧ x'=2",
			expression: Conjunction(
				NewBinary(NewVariable(y), Equal, NewInteger(1)),
				NewAssignment(NewVariable(x), NewInteger(2)),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations state.Valuations[*z3.AST]) {
				if valuation, exists := valuations.Value(x); exists {
					if !solver.Proven(z3.Eq(valuation, context.NewInt(2, context.IntegerSort()))) {
						t.Errorf("Expected x to be 2 but was %s", valuation.String())
					}
				} else {
					t.Errorf("Expected x in valuations")
				}
			},
		},
		{
			name: "x'=2 ∧ x=2 ∧ x'=3",
			expression: Conjunction(
				NewAssignment(NewVariable(x), NewInteger(2)),
				NewBinary(NewVariable(x), Equal, NewInteger(2)),
				NewAssignment(NewVariable(x), NewInteger(3)),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations state.Valuations[*z3.AST]) {
				if valuation, exists := valuations.Value(x); exists {
					if !solver.Proven(z3.Eq(valuation, context.NewInt(3, context.IntegerSort()))) {
						t.Errorf("Expected x to be 3 but was %s", valuation.String())
					}
				} else {
					t.Errorf("Expected x in valuations")
				}
			},
		},
		{
			name: "x=3 ∨ (x=1 ∧ x=2)",
			expression: Disjunction(
				NewBinary(NewVariable(x), Equal, NewInteger(3)),
				Conjunction(
					NewBinary(NewVariable(x), Equal, NewInteger(1)),
					NewBinary(NewVariable(x), Equal, NewInteger(2)),
				),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations state.Valuations[*z3.AST]) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
		{
			name: "x=1 ∨ x=2",
			expression: Disjunction(
				NewBinary(NewVariable(x), Equal, NewInteger(1)),
				NewBinary(NewVariable(x), Equal, NewInteger(2)),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations state.Valuations[*z3.AST]) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valuations := state.NewValuationsMap[*z3.AST]()
			valuations.Assign(y, context.NewInt(1, context.IntegerSort()))
			interpreter := NewSymbolicInterpreter(context, variables, valuations)

			actual := interpreter.Interpret(tt.expression)
			if !solver.Proven(z3.Eq(actual, tt.expected)) {
				assert.Equal(t, tt.expected.String(), actual.String())
			}

			if tt.valuations != nil {
				tt.valuations(t, valuations)
			}
		})
	}
}
