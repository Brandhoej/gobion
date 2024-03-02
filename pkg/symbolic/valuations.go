package symbolic

type Valuations[T any] interface {
	Load(symbol Symbol) (T, bool)
	Store(symbol Symbol, value T)
}

type ValuationsMap[T any] struct {
	valuations map[Symbol]T
}

func NewEnvironmentMap[T any]() *ValuationsMap[T] {
	return &ValuationsMap[T]{
		valuations: map[Symbol]T{},
	}
}

func (environment *ValuationsMap[T]) Load(symbol Symbol) (value T, exists bool) {
	value, exists = environment.valuations[symbol]
	return
}

func (environment *ValuationsMap[T]) Store(symbol Symbol, value T) {
	environment.valuations[symbol] = value
}
