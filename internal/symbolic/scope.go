package symbolic

type Scope[T any] interface {
	Symbols
	Environment[T]
}

type GoScope[T any] struct {
	parent *GoScope[T]
	symbols Symbols
	environment Environment[T]
}

func (scope *GoScope[T]) Child() *GoScope[T] {
	return &GoScope[T]{
		parent: scope,
		symbols: NewSymbolsMap(),
	}
}

func (scope *GoScope[T]) Insert(identifier string) Symbol {
	return scope.symbols.Insert(identifier)
}

func (scope *GoScope[T]) Contains(symbol Symbol) bool {
	if scope.symbols.Contains(symbol) {
		return true
	}

	if scope.parent != nil {
		return scope.parent.Contains(symbol)
	}

	return false
}

func (scope *GoScope[T]) Load(symbol Symbol) (T, bool) {
	if value, exists := scope.environment.Load(symbol); exists {
		return value, exists
	}

	if scope.parent != nil {
		return scope.parent.environment.Load(symbol)
	}

	var zero T
	return zero, false
}

func (scope *GoScope[T]) Store(symbol Symbol, value T) {
	scope.Store(symbol, value)
}
