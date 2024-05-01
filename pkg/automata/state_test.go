package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/graph"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_SubsetOf(t *testing.T) {
	symbols := symbols.NewSymbolsMap[string](symbols.NewSymbolsFactory())
	x := language.NewVariable(symbols.Insert("x"))
	valuations := language.NewValuationsMap()
	context := z3.NewContext(z3.NewConfig())
	variables := language.NewVariablesMap()
	variables.Declare(symbols.Insert("x"), language.IntegerSort)
	solver := NewInterpreter(context, variables)
	tests := []struct {
		name     string
		lhs, rhs State
		expected bool
	}{
		{
			name:     "(0, True) ⊆ (0, True)",
			lhs:      NewState(graph.Key(0), valuations, language.NewTrue()),
			rhs:      NewState(graph.Key(0), valuations, language.NewTrue()),
			expected: true,
		},
		{
			name:     "(1, True) ⊈ (0, True)",
			lhs:      NewState(graph.Key(1), valuations, language.NewTrue()),
			rhs:      NewState(graph.Key(0), valuations, language.NewTrue()),
			expected: false,
		},
		{
			name:     "(0, False) ⊈ (0, True)",
			lhs:      NewState(graph.Key(0), valuations, language.NewFalse()),
			rhs:      NewState(graph.Key(0), valuations, language.NewTrue()),
			expected: true,
		},
		{
			name: "(0, {x=0}) ⊆ (0, True)",
			lhs: NewState(
				graph.Key(0),
				valuations,
				language.NewBinary(
					x, language.Equal, language.NewInteger(0),
				),
			),
			rhs:      NewState(graph.Key(0), valuations, language.NewTrue()),
			expected: true,
		},
		{
			name: "(0, True) ⊈ (0, {x=0})",
			lhs: NewState(
				graph.Key(0),
				valuations,
				language.NewBinary(
					x, language.Equal, language.NewInteger(0),
				),
			),
			rhs:      NewState(graph.Key(0), valuations, language.NewTrue()),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.lhs.SubsetOf(tt.rhs, solver)
			assert.Equal(t, tt.expected, actual, tt.name)
		})
	}
}
