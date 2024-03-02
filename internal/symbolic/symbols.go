package symbolic

type Symbol int32

type Symbols interface {
	Insert(identifier string) Symbol
	Contains(symbol Symbol) bool
}

type SymbolsMap struct {
	identifiers map[string]Symbol
	symbols map[Symbol]string
}

func NewSymbolsMap() *SymbolsMap {
	return &SymbolsMap{
		identifiers: map[string]Symbol{},
		symbols: map[Symbol]string{},
	}
}

func (symbols *SymbolsMap) Next() Symbol {
	return Symbol(len(symbols.identifiers))
}

func (symbols *SymbolsMap) Insert(identifier string) Symbol {
	if symbol, exists := symbols.identifiers[identifier]; exists {
		return symbol
	}
	return symbols.Next()
}

func (symbols *SymbolsMap) Contains(symbol Symbol) bool {
	_, exists := symbols.symbols[symbol]
	return exists
}