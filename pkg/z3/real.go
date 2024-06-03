package z3

type Real struct {
	numerator, denominator int
}

func NewReal(numerator, denominator int) Real {
	return Real{
		numerator:  numerator,
		denominator: denominator,
	}
}