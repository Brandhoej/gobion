package expressions

type Assertions[T any] interface {
	Constrain(constraint T)
	All(yield func(assertion T) bool) bool
}
