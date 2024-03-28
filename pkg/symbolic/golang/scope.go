package golang

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbolic"
)

type GoScope struct {
	id         int
	parent     *GoScope
	symbols    symbolic.Symbols
	valuations symbolic.Valuations
	variables  symbolic.Variables
	solver     *z3.Solver
}

func NewGoGlobalScope(id int, solver *z3.Solver) *GoScope {
	return &GoScope{
		id:         id,
		parent:     nil,
		symbols:    symbolic.NewSymbolsMap(),
		valuations: symbolic.NewEnvironmentMap(),
		variables:  symbolic.NewVariablesMap(),
		solver:     solver,
	}
}

func (scope *GoScope) Parent() *GoScope {
	return scope.parent
}

func (scope *GoScope) Child(id int) *GoScope {
	return &GoScope{
		id:         id,
		parent:     scope,
		symbols:    symbolic.NewSymbolsMap(),
		valuations: symbolic.NewEnvironmentMap(),
		variables:  symbolic.NewVariablesMap(),
		solver:     scope.CopyState(),
	}
}

func (scope *GoScope) State() *z3.Solver {
	return scope.solver
}

func (scope *GoScope) CopyState() *z3.Solver {
	solver := scope.solver.Context().NewSolver()
	var recursive func(scope *GoScope)
	recursive = func(scope *GoScope) {
		if scope == nil {
			return
		}

		// Provided the symbolic states we have and our path constraint.
		// Is there always a solution.
		for _, symbol := range scope.symbols.Symbols() {
			variable, _ := scope.VariableFor(symbol)
			valuation, _ := scope.ValuationFor(symbol)
			solver.Assert(z3.Eq(variable, valuation))
		}

		recursive(scope.parent)
	}
	recursive(scope)

	return solver
}

func (scope *GoScope) Symbol(identifier string) (symbolic.Symbol, *GoScope) {
	symbol, exists := scope.symbols.Lookup(identifier)
	if exists {
		return symbol, scope
	}

	// The identifier is not present in this scope but might be in the parent.
	if scope.parent != nil {
		return scope.parent.Symbol(identifier)
	}

	// The identifier could not be found in this scope and it has no parent.
	return symbol, nil
}

// Returns true if the identifier is declared in the scope and not one of it ancestors.
func (scope *GoScope) IsLocal(identifier string) bool {
	_, exists := scope.symbols.Lookup(identifier)
	return exists
}

// Returns true if the identifier is local to this scope or it is accessible in its parent.
func (scope *GoScope) IsAccessible(identifier string) bool {
	if scope.IsLocal(identifier) {
		return true
	}

	if scope.parent != nil {
		return scope.parent.IsAccessible(identifier)
	}

	return false
}

// Returns true if the scope is the global scope of the program (i.e., it does not have a parent scope).
func (scope *GoScope) IsGlobal() bool {
	return scope.parent == nil
}

// Declares a variable with a zero-valuation.
func (scope *GoScope) Define(identifier string, sort *z3.Sort) (variable *z3.AST) {
	var valuation *z3.AST
	switch sort.Kind() {
	case z3.KindInt:
		valuation = sort.Context().NewInt(0, sort)
	case z3.KindBoolean:
		valuation = sort.Context().NewFalse()
	}

	return scope.Declare(identifier, valuation)
}

// Declares a variable with the valuation.
func (scope *GoScope) Declare(identifier string, valuation *z3.AST) *z3.AST {
	symbol := scope.symbols.Insert(identifier)
	scopedIdentifier := fmt.Sprintf("%s_%d", identifier, scope.id)
	variable := scope.variables.Declare(symbol, scopedIdentifier, valuation.Sort())
	scope.valuations.Store(symbol, valuation)
	scope.solver.Assert(z3.Eq(variable, valuation))
	return variable
}

// Assigned a new valuation to a varaible.
// This updates the variable stored for the identifier to ensure Single Static Assignment (SSA).
func (scope *GoScope) Assign(identifier string, valuation *z3.AST) *z3.AST {
	symbol, scope := scope.Symbol(identifier)
	if scope == nil {
		panic("Cannot assign to undeclared identifier")
	}
	return scope.AssignTo(symbol, valuation)
}

func (scope *GoScope) AssignTo(symbol symbolic.Symbol, valuation *z3.AST) *z3.AST {
	variable, _ := scope.variables.Advance(symbol)
	scope.valuations.Store(symbol, valuation)
	scope.solver.Assert(z3.Eq(variable, valuation))
	return variable
}

// Looks up a valuation. If the variable is bound that the valuation is the bound variable itself.
// However, if the variable is not bound then the valuation is the reduced variables.
func (scope *GoScope) Valuation(identifier string) (valuation *z3.AST, exists bool) {
	symbol, scope := scope.Symbol(identifier)
	if scope == nil {
		return
	}

	return scope.ValuationFor(symbol)
}

func (scope *GoScope) ValuationFor(symbol symbolic.Symbol) (valuation *z3.AST, exists bool) {
	valuation, exists = scope.valuations.Load(symbol)
	if !exists && scope.parent != nil {
		valuation, exists = scope.parent.ValuationFor(symbol)
	}
	return
}

// Looks up a variable and returns it.
func (scope *GoScope) Variable(identifier string) (variable *z3.AST, exists bool) {
	symbol, scope := scope.Symbol(identifier)
	if scope == nil {
		return
	}

	return scope.VariableFor(symbol)
}

// Looks up a variable and returns it.
func (scope *GoScope) VariableFor(symbol symbolic.Symbol) (variable *z3.AST, exists bool) {
	variable, exists = scope.variables.Variable(symbol)
	if !exists && scope.parent != nil {
		variable, exists = scope.parent.VariableFor(symbol)
	}
	return
}

func (scope *GoScope) Traverse(observer func(scope *GoScope)) {
	var recursive func(*GoScope)
	recursive = func(scope *GoScope) {
		if scope == nil {
			return
		}
		observer(scope)
		recursive(scope.parent)
	}
	recursive(scope)
}

// Returns all the local variables' identifiers.
func (scope *GoScope) Identifiers() []string {
	identifiers := scope.symbols.Identifiers()
	if scope.parent != nil {
		identifiers = append(identifiers, scope.parent.Identifiers()...)
	}
	return identifiers
}

func (scope *GoScope) string(builder *strings.Builder) {
	if scope.parent != nil {
		scope.parent.string(builder)
	}

	for _, identifier := range scope.symbols.Identifiers() {
		symbol, _ := scope.symbols.Lookup(identifier)
		if value, exits := scope.valuations.Load(symbol); exits {
			builder.WriteString(fmt.Sprintf("%v=%v ", identifier, value))
		} else {
			builder.WriteString(fmt.Sprintf("%v=?? ", identifier))
		}
	}
}

func (scope *GoScope) String() string {
	var builder strings.Builder
	scope.string(&builder)
	return builder.String()
}
