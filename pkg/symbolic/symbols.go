package symbolic

type Symbol int32

type Symbols interface {
	Insert(identifier string) Symbol
	Lookup(identifier string) (Symbol, bool)
	Identifiers() []string
	Symbols() []Symbol
}

type SymbolsMap struct {
	identifiers map[string]Symbol
	symbols     map[Symbol]string
}

func NewSymbolsMap() *SymbolsMap {
	return &SymbolsMap{
		identifiers: map[string]Symbol{},
		symbols:     map[Symbol]string{},
	}
}

func (mapping *SymbolsMap) Next() Symbol {
	return Symbol(len(mapping.identifiers))
}

func (mapping *SymbolsMap) Insert(identifier string) Symbol {
	if symbol, exists := mapping.identifiers[identifier]; exists {
		return symbol
	}
	symbol := mapping.Next()
	mapping.identifiers[identifier] = symbol
	return symbol
}

func (mapping *SymbolsMap) Lookup(identifier string) (Symbol, bool) {
	symbol, exists := mapping.identifiers[identifier]
	return symbol, exists
}

func (mapping *SymbolsMap) Identifiers() []string {
	identifiers := make([]string, 0, len(mapping.identifiers))
	for identifier := range mapping.identifiers {
		identifiers = append(identifiers, identifier)
	}
	return identifiers
}

func (mapping *SymbolsMap) Symbols() []Symbol {
	symbols := make([]Symbol, 0, len(mapping.identifiers))
	for _, symbol := range mapping.identifiers {
		symbols = append(symbols, symbol)
	}
	return symbols
}
