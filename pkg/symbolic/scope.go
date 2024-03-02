package symbolic

import (
	"fmt"
	"strings"
)

type Scope[T any] interface {
	Declare(identifier string, variable, value T)
	Symbol(identifier string) (symbol Symbol, exists bool)
	IsLocal(identifier string) bool
	IsGlobal() bool
	Assign(symbol Symbol, value T)
	Variable(symbol Symbol) (variable T, exists bool)
	Valuation(symbol Symbol) (valuation T, exists bool)
}

type GoScope[T any] struct {
	parent      *GoScope[T]
	symbols     Symbols
	valuations  Valuations[T]
	variables   Variables[T]
}

func NewGoGlobalScope[T any]() *GoScope[T] {
	return &GoScope[T]{
		parent:      nil,
		symbols:     NewSymbolsMap(),
		valuations: NewEnvironmentMap[T](),
		variables: NewVariablesMap[T](),
	}
}

func (scope *GoScope[T]) Parent() *GoScope[T] {
	return scope.parent
}

func (scope *GoScope[T]) Child() *GoScope[T] {
	return &GoScope[T]{
		parent:      scope,
		symbols:     NewSymbolsMap(),
		valuations: NewEnvironmentMap[T](),
		variables: NewVariablesMap[T](),
	}
}

func (scope *GoScope[T]) Declare(identifier string, variable, value T) {
	symbol := scope.symbols.Insert(identifier)
	scope.variables.Declare(symbol, variable)
	scope.valuations.Store(symbol, value)
}

func (scope *GoScope[T]) Symbol(identifier string) (symbol Symbol, exists bool) {
	symbol, exists = scope.symbols.Lookup(identifier)
	if !exists && scope.parent != nil {
		symbol, exists = scope.parent.Symbol(identifier)
	}
	return symbol, exists
}

func (scope *GoScope[T]) IsLocal(identifier string) bool {
	_, exists := scope.symbols.Lookup(identifier)
	return exists
}

func (scope *GoScope[T]) IsGlobal() bool {
	return scope.parent == nil
}

func (scope *GoScope[T]) AssignTo(identifier string, value T) {
	symbol, exists := scope.Symbol(identifier)
	if !exists {
		panic("Cannot assign to undeclared identifier")
	}
	scope.Assign(symbol, value)
}

func (scope *GoScope[T]) Assign(symbol Symbol, value T) {
	scope.valuations.Store(symbol, value)
}

func (scope *GoScope[T]) Variable(symbol Symbol) (variable T, exists bool) {
	variable, exists = scope.variables.Variable(symbol)
	if !exists && scope.parent != nil {
		variable, exists = scope.parent.Variable(symbol)
	}
	return
}

func (scope *GoScope[T]) Valuation(symbol Symbol) (valuation T, exists bool) {
	valuation, exists = scope.valuations.Load(symbol)
	if !exists && scope.parent != nil {
		valuation, exists = scope.parent.Valuation(symbol)
	}
	return
}

func (scope *GoScope[T]) Identifiers() []string {
	identifiers := scope.symbols.Identifiers()
	if scope.parent != nil {
		identifiers = append(identifiers, scope.parent.Identifiers()...)
	}
	return identifiers
}

func (scope *GoScope[T]) Symbols() []Symbol {
	identifiers := scope.symbols.Symbols()
	if scope.parent != nil {
		identifiers = append(identifiers, scope.parent.Symbols()...)
	}
	return identifiers
}


func (scope *GoScope[T]) string(builder *strings.Builder) {
	for _, identifier := range scope.symbols.Identifiers() {
		symbol, _ := scope.symbols.Lookup(identifier)
		if value, exits := scope.valuations.Load(symbol); exits {
			builder.WriteString(fmt.Sprintf("%v=%v ", identifier, value))
		} else {
			builder.WriteString(fmt.Sprintf("%v=?? ", identifier))
		}
	}
}

func (scope *GoScope[T]) String() string {
	var builder strings.Builder

	if scope.parent != nil {
		scope.string(&builder)
	}
	scope.string(&builder)

	return builder.String()
}