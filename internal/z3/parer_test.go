package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	// Arrange
	config := NewConfig()
	context := NewContext(config)

	// Act
	asts := context.Parse(`
	(declare-const x Int)
	(assert (= x 10))
	(define-fun max_integ ((x Int) (y Int)) Int 
		(ite (< x y) y x))
	(define-fun min_integ ((x Int) (y Int)) Int 
		(ite (< x y) x y))
		(assert (= (min_integ 10 0) 10))
		(assert (= (max_integ 10 0) 0))
	`)

	// Assert
	assert.Equal(t, uint(3), asts.Length())
	assert.Equal(t, "(= x 10)", asts.Get(0).String())
	assert.Equal(t, "(= (ite (< 10 0) 10 0) 10)", asts.Get(1).String())
	assert.Equal(t, "(= (ite (< 10 0) 0 10) 0)", asts.Get(2).String())
}
