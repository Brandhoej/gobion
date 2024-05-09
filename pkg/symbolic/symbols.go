package symbolic

import "github.com/Brandhoej/gobion/internal/z3"

type Symbol int32

func (symbol Symbol) Z3(context *z3.Context) z3.Symbol {
	return context.NewIntegerSymbol(int(symbol))
}

type Symbols interface {
	Insert(identifier string) (symbol Symbol)
	Lookup(identifier string) (symbol Symbol, exists bool)
	Identifier(symbol Symbol) (identifier string, exists bool)
}

type SymbolsMap struct {
	factory     *SymbolsFactory
	symbols     map[Symbol]string
	identifiers map[string]Symbol
}

func NewSymbolsMap(factory *SymbolsFactory) *SymbolsMap {
	return &SymbolsMap{
		factory:     factory,
		symbols:     map[Symbol]string{},
		identifiers: map[string]Symbol{},
	}
}

func (mapping *SymbolsMap) Insert(identifier string) (symbol Symbol) {
	if symbol, exists := mapping.identifiers[identifier]; exists {
		return symbol
	}
	symbol = mapping.factory.Next()
	mapping.symbols[symbol] = identifier
	mapping.identifiers[identifier] = symbol
	return symbol
}

func (mapping *SymbolsMap) Lookup(identifier string) (symbol Symbol, exists bool) {
	symbol, exists = mapping.identifiers[identifier]
	return symbol, exists
}

func (mapping *SymbolsMap) Identifier(symbol Symbol) (identifier string, exists bool) {
	identifier, exists = mapping.symbols[symbol]
	return identifier, exists
}
