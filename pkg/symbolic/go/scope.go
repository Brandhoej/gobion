package symbolic

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbolic"
)

type GoScope struct {
	depth      int
	parent     *GoScope
	symbols    symbolic.Symbols
	valuations symbolic.Valuations
	variables  symbolic.Variables
}

func NewGoGlobalScope() *GoScope {
	return &GoScope{
		depth:      0,
		parent:     nil,
		symbols:    symbolic.NewSymbolsMap(),
		valuations: symbolic.NewEnvironmentMap(),
		variables:  symbolic.NewVariablesMap(),
	}
}

func (scope *GoScope) Parent() *GoScope {
	return scope.parent
}

func (scope *GoScope) Child() *GoScope {
	return &GoScope{
		depth:      scope.depth + 1,
		parent:     scope,
		symbols:    symbolic.NewSymbolsMap(),
		valuations: symbolic.NewEnvironmentMap(),
		variables:  symbolic.NewVariablesMap(),
	}
}

func (scope *GoScope) symbol(identifier string) (symbolic.Symbol, *GoScope) {
	symbol, exists := scope.symbols.Lookup(identifier)
	if exists {
		return symbol, scope
	}

	// The identifier is not present in this scope but might be in the parent.
	if scope.parent != nil {
		return scope.parent.symbol(identifier)
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
	prefixedIdentifier := fmt.Sprintf("%s_%d", identifier, scope.depth)
	variable := scope.variables.Declare(symbol, prefixedIdentifier, valuation.Sort())
	scope.valuations.Store(symbol, valuation)
	return variable
}

// Assigned a new valuation to a varaible.
// This updates the variable stored for the identifier to ensure Single Static Assignment (SSA).
func (scope *GoScope) Assign(identifier string, valuation *z3.AST) *z3.AST {
	symbol, scope := scope.symbol(identifier)
	if scope == nil {
		panic("Cannot assign to undeclared identifier")
	}
	variable, _ := scope.variables.Advance(symbol)
	scope.valuations.Store(symbol, valuation)
	return variable
}

// Looks up a variable and returns it.
func (scope *GoScope) Variable(identifier string) (variable *z3.AST, exists bool) {
	if symbol, hasSymbol := scope.symbols.Lookup(identifier); hasSymbol {
		variable, exists = scope.variables.Variable(symbol)
		return
	}
	if !exists && scope.parent != nil {
		variable, exists = scope.parent.Variable(identifier)
	}
	return
}

// Looks up a valuation. If the variable is bound that the valuation is the bound variable itself.
// However, if the variable is not bound then the valuation is the reduced variables.
func (scope *GoScope) Valuation(identifier string) (valuation *z3.AST, exists bool) {
	symbol, scope := scope.symbol(identifier)
	if scope == nil {
		return
	}

	valuation, exists = scope.valuations.Load(symbol)
	if !exists && scope.parent != nil {
		valuation, exists = scope.parent.Valuation(identifier)
	}
	return
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
