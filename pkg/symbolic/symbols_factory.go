package symbolic

type SymbolsFactory struct {
	counter int
}

func NewSymbolsFactory() *SymbolsFactory {
	return &SymbolsFactory{
		counter: 0,
	}
}

func (factory *SymbolsFactory) Next() Symbol {
	id := factory.counter
	factory.counter += 1
	return Symbol(id)
}