package symbolic

type Environment[T any] interface {
	Load(symbol Symbol) (T, bool)
	Store(symbol Symbol, value T)
}

type EnvironmentMap[T any] struct {
	variables map[Symbol]T
}

func (environment EnvironmentMap[T]) Load(symbol Symbol) (value T, exists bool) {
	value, exists = environment.variables[symbol]
	return
}

func (environment EnvironmentMap[T]) Store(symbol Symbol, value T) {
	environment.variables[symbol] = value
}
