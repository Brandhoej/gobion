package symbols

type Symbol int

type Store[T any] interface {
	Insert(item T) (symbol Symbol)
	Lookup(item T) (symbol Symbol, exists bool)
	Item(symbol Symbol) (identifier T, exists bool)
}

type SymbolsMap[T comparable] struct {
	factory     *SymbolsFactory
	symbols     map[Symbol]T
	identifiers map[T]Symbol
}

func NewSymbolsMap[T comparable](factory *SymbolsFactory) *SymbolsMap[T] {
	return &SymbolsMap[T]{
		factory:     factory,
		symbols:     map[Symbol]T{},
		identifiers: map[T]Symbol{},
	}
}

func (mapping *SymbolsMap[T]) Insert(item T) (symbol Symbol) {
	if symbol, exists := mapping.identifiers[item]; exists {
		return symbol
	}
	symbol = mapping.factory.Next()
	mapping.symbols[symbol] = item
	mapping.identifiers[item] = symbol
	return symbol
}

func (mapping *SymbolsMap[T]) Lookup(item T) (symbol Symbol, exists bool) {
	symbol, exists = mapping.identifiers[item]
	return symbol, exists
}

func (mapping *SymbolsMap[T]) Item(symbol Symbol) (item T, exists bool) {
	item, exists = mapping.symbols[symbol]
	return item, exists
}
