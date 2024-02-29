package z3

// Prevents user code from directly comparing values for equality.
//
// It uses an anonymous field _ with type [0]func() that a zero-sized array of functions.
// It does not have any practical purpose other than to ensure that instances
// of the noEq struct have a unique type distinct from any other struct type.
type noEq struct {
	_ [0]func()
}
