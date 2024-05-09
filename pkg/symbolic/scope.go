package symbolic

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/structures"
)

type Scope interface {
	Bind(identifier string) (symbol Symbol)
	Lookup(identifier string) (symbol Symbol, encasing Scope)
	Identifier(symbol Symbol) (identifier string, exists bool)

	Define(symbol Symbol, sort *z3.Sort) (variable *z3.AST)
	Declare(symbol Symbol, valuation *z3.AST) (variable *z3.AST)
	Assign(symbol Symbol, valuation *z3.AST)

	Variable(symbol Symbol) (variable *z3.AST)
	Variables(function func(symbol Symbol, scope Scope))
	Valuation(symbol Symbol) (valuation *z3.AST)

	Branch() (branch Scope)
	Solver() *z3.Solver
}

type LexicalScope struct {
	context        *z3.Context
	parent         *LexicalScope
	symbolsFactory *SymbolsFactory
	symbols        Symbols
	variables      Variables
	valuations     Valuations
}

func NewLexicalScope(context *z3.Context) *LexicalScope {
	symbolsFactory := NewSymbolsFactory()
	return &LexicalScope{
		parent:         nil,
		context:        context,
		symbolsFactory: symbolsFactory,
		symbols:        NewSymbolsMap(symbolsFactory),
		variables:      NewVariablesMap(),
		valuations:     NewValuationsMap(),
	}
}

func (scope *LexicalScope) Branch() (branch Scope) {
	return &LexicalScope{
		parent:         scope,
		context:        scope.context,
		symbolsFactory: scope.symbolsFactory,
		symbols:        NewSymbolsMap(scope.symbolsFactory),
		variables:      NewVariablesMap(),
		valuations:     NewValuationsMap(),
	}
}

func (scope *LexicalScope) Bind(identifier string) (symbol Symbol) {
	return scope.symbols.Insert(identifier)
}

func (scope *LexicalScope) Lookup(identifier string) (delcaration Symbol, encasing Scope) {
	scope.traverse(func(scope *LexicalScope) bool {
		if symbol, exists := scope.symbols.Lookup(identifier); exists {
			delcaration = symbol
			encasing = scope
			return false
		}
		return true
	})
	return delcaration, encasing
}

func (scope *LexicalScope) Identifier(symbol Symbol) (identifier string, exists bool) {
	scope.traverse(func(scope *LexicalScope) bool {
		identifier, exists = scope.Identifier(symbol)
		return !exists
	})
	return identifier, exists
}

func (scope *LexicalScope) Define(symbol Symbol, sort *z3.Sort) (variable *z3.AST) {
	return scope.Declare(symbol, sort.Zero())
}

func (scope *LexicalScope) Declare(symbol Symbol, valuation *z3.AST) *z3.AST {
	// We cannot re-declare the same symbol in the scope.
	if variable := scope.variables.Get(symbol); variable != nil {
		return nil
	}

	// Assign a value to the new declaration in this scope.
	scope.valuations.Store(symbol, valuation)
	return scope.variables.Declare(symbol, valuation.Sort())
}

func (scope *LexicalScope) Assign(symbol Symbol, valuation *z3.AST) {
	scope.traverse(func(scope *LexicalScope) bool {
		variable := scope.variables.Get(symbol)
		if variable != nil {
			scope.valuations.Store(symbol, valuation)
		}
		return variable == nil
	})
}

func (scope *LexicalScope) Variable(symbol Symbol) (variable *z3.AST) {
	scope.traverse(func(scope *LexicalScope) bool {
		variable = scope.variables.Get(symbol)
		return variable == nil
	})

	return variable
}

func (scope *LexicalScope) Valuation(symbol Symbol) (valuation *z3.AST) {
	scope.traverse(func(scope *LexicalScope) bool {
		valuation = scope.valuations.Load(symbol)
		return valuation == nil
	})

	return valuation
}

func (scope *LexicalScope) Variables(function func(symbol Symbol, scope Scope)) {
	scope.traverse(func(scope *LexicalScope) bool {
		scope.variables.Symbols(func(symbol Symbol) {
			function(symbol, scope)
		})
		return true
	})
}

func (scope *LexicalScope) traverse(function func(scope *LexicalScope) bool) {
	structures.BottomUp(
		scope,
		func(current *LexicalScope) (proceeding *LexicalScope, exists bool) {
			return current.parent, current.parent != nil
		},
		func(current *LexicalScope) bool {
			return function(current)
		},
	)
}

func (scope *LexicalScope) Solver() *z3.Solver {
	solver := scope.context.NewSolver()
	scope.Variables(func(symbol Symbol, current Scope) {
		variable := current.Variable(symbol)
		valuation := scope.Valuation(symbol)
		equality := z3.Eq(variable, valuation)
		solver.Assert(equality)
	})
	return solver
}
