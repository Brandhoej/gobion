package symbolic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbolic"
)

type GoScope struct {
	parent     *GoScope
	symbols    symbolic.Symbols
	valuations symbolic.Valuations
	variables  symbolic.Variables
}

func NewGoGlobalScope() *GoScope {
	return &GoScope{
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
		parent:     scope,
		symbols:    symbolic.NewSymbolsMap(),
		valuations: symbolic.NewEnvironmentMap(),
		variables:  symbolic.NewVariablesMap(),
	}
}

func (scope *GoScope) IsLocal(identifier string) bool {
	_, exists := scope.symbols.Lookup(identifier)
	return exists
}

func (scope *GoScope) IsGlobal() bool {
	return scope.parent == nil
}

func (scope *GoScope) Symbol(identifier string) (symbol symbolic.Symbol, exists bool) {
	symbol, exists = scope.symbols.Lookup(identifier)
	if !exists && scope.parent != nil {
		symbol, exists = scope.parent.Symbol(identifier)
	}
	return symbol, exists
}

func (scope *GoScope) Bind(index uint, sort *z3.Sort) (variable *z3.AST) {
	var valuation *z3.AST
	switch sort.Kind() {
	case z3.KindInt:
		valuation = sort.Context().NewInt(0, sort)
	case z3.KindBoolean:
		valuation = sort.Context().NewFalse()
	}
	symbol := scope.symbols.Insert(strconv.Itoa(int(index)))
	variable = scope.variables.Bind(index, sort)
	scope.valuations.Store(symbol, valuation)
	return
}

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

func (scope *GoScope) Declare(identifier string, valuation *z3.AST) *z3.AST {
	symbol := scope.symbols.Insert(identifier)
	variable := scope.variables.Declare(symbol, identifier, valuation.Sort())
	scope.valuations.Store(symbol, valuation)
	return variable
}

func (scope *GoScope) Assign(identifier string, valuation *z3.AST) *z3.AST {
	symbol, exists := scope.Symbol(identifier)
	if !exists {
		panic("Cannot assign to undeclared identifier")
	}
	variable, _ := scope.variables.Advance(symbol)
	scope.valuations.Store(symbol, valuation)
	return variable
}

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

func (scope *GoScope) Valuation(identifier string) (valuation *z3.AST, exists bool) {
	symbol, exists := scope.Symbol(identifier)
	if !exists {
		return
	}

	valuation, exists = scope.valuations.Load(symbol)
	if !exists && scope.parent != nil {
		valuation, exists = scope.parent.Valuation(identifier)
	}
	return
}

func (scope *GoScope) Identifiers() []string {
	identifiers := scope.symbols.Identifiers()
	if scope.parent != nil {
		identifiers = append(identifiers, scope.parent.Identifiers()...)
	}
	return identifiers
}

func (scope *GoScope) Symbols() []symbolic.Symbol {
	identifiers := scope.symbols.Symbols()
	if scope.parent != nil {
		identifiers = append(identifiers, scope.parent.Symbols()...)
	}
	return identifiers
}

func (scope *GoScope) string(builder *strings.Builder) {
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

	if scope.parent != nil {
		scope.string(&builder)
	}
	scope.string(&builder)

	return builder.String()
}
