package symbolic

type Variables[T any] interface {
	Declare(symbol Symbol, variable T)
	Variable(symbol Symbol) (variable T, exists bool)
}

type VariablesMap[T any] struct {
	variables map[Symbol]T
}

func NewVariablesMap[T any]() *VariablesMap[T] {
	return &VariablesMap[T]{
		variables: map[Symbol]T{},
	}
}

func (environment *VariablesMap[T]) Declare(symbol Symbol, variable T) {
	environment.variables[symbol] = variable
}

func (environment *VariablesMap[T]) Variable(symbol Symbol) (variable T, exists bool) {
	variable, exists = environment.variables[symbol]
	return variable, exists
}