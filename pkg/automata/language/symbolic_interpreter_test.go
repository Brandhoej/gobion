package language

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
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

	variables := NewVariablesMap()
	variables.Declare(x, IntegerSort)
	variables.Declare(y, IntegerSort)

	tests := []struct {
		name       string
		expression Expression
		expected   *z3.AST
		valuations func(t *testing.T, valuations Valuations)
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
			name: "x+y",
			expression: NewBinary(
				NewVariable(x),
				Addition,
				NewVariable(y),
			),
			expected: z3.Add(
				xConst, context.NewInt(1, context.IntegerSort()),
			),
		},
		{
			name: "x>0 ? 1 : 2",
			expression: NewIfThenElse(
				NewBinary(
					NewVariable(x),
					GreaterThan,
					NewInteger(0),
				),
				NewInteger(1),
				NewInteger(2),
			),
			expected: context.NewInt(2, context.IntegerSort()),
		},
		{
			name: "x=2",
			expression: NewBinary(
				NewVariable(x),
				Equal,
				NewInteger(2),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations Valuations) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
		{
			name: "y=3 ∧ {x'=2; true}",
			expression: Conjunction(
				NewBinary(
					NewVariable(y),
					Equal,
					NewInteger(3),
				),
				NewBlockExpression(
					NewTrue(),
					NewAssignment(
						NewVariable(x),
						NewInteger(2),
					),
				),
			),
			expected: context.NewFalse(),
			valuations: func(t *testing.T, valuations Valuations) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
		{
			name: "y=1 ∨ {x'=2; true}",
			expression: Disjunction(
				NewBinary(
					NewVariable(y),
					Equal,
					NewInteger(1),
				),
				NewBlockExpression(
					NewTrue(),
					NewAssignment(
						NewVariable(x),
						NewInteger(2),
					),
				),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations Valuations) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
		{
			name: "y=1 ∧ {x'=2; true}",
			expression: Conjunction(
				NewBinary(
					NewVariable(y),
					Equal,
					NewInteger(1),
				),
				NewBlockExpression(
					NewTrue(),
					NewAssignment(
						NewVariable(x),
						NewInteger(2),
					),
				),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations Valuations) {
				if valuation, exists := valuations.Value(x); exists {
					if integer, ok := valuation.(Integer); !(ok && integer.value == 2) {
						t.Errorf("Expected x to be 2 but was %v", integer.value)
					}
				} else {
					t.Errorf("Expected x in valuations")
				}
			},
		},
		{
			name: "{x'=2; true} ∧ x=2 ∧ {x'=3; true}",
			expression: Conjunction(
				NewBlockExpression(
					NewTrue(),
					NewAssignment(
						NewVariable(x),
						NewInteger(2),
					),
				),
				NewBinary(
					NewVariable(x),
					Equal,
					NewInteger(2),
				),
				NewBlockExpression(
					NewTrue(),
					NewAssignment(
						NewVariable(x),
						NewInteger(3),
					),
				),
			),
			expected: context.NewTrue(),
			valuations: func(t *testing.T, valuations Valuations) {
				if valuation, exists := valuations.Value(x); exists {
					if integer, ok := valuation.(Integer); !(ok && integer.value == 3) {
						t.Errorf("Expected x to be 3 but was %v", integer.value)
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
			valuations: func(t *testing.T, valuations Valuations) {
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
			valuations: func(t *testing.T, valuations Valuations) {
				if _, exists := valuations.Value(x); exists {
					t.Errorf("Did not expect a valuation of x")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valuations := NewValuationsMap()
			valuations.Assign(y, NewInteger(1))
			interpreter := NewSymbolicInterpreter(context, variables, valuations)

			actual := interpreter.Expression(tt.expression)
			if !solver.Proven(z3.Eq(actual, tt.expected)) {
				assert.Equal(t, tt.expected.String(), actual.String())
			}

			if tt.valuations != nil {
				tt.valuations(t, valuations)
			}
		})
	}
}
