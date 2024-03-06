package z3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	var context *Context

	boolV := func(identifier string) *AST {
		return context.NewConstant(WithName(identifier), context.BooleanSort())
	}
	realV := func(identifier string) *AST {
		return context.NewConstant(WithName(identifier), context.RealSort())
	}
	intV := func(identifier string) *AST {
		return context.NewConstant(WithName(identifier), context.IntegerSort())
	}
	intC := func(value int) *AST {
		return context.NewInt(value, context.IntegerSort())
	}

	tests := []struct {
		name      string
		operation func() *AST
		text      string
		kind      Kind
	}{
		{
			name: "Logical conjunction",
			operation: func() *AST {
				return And(boolV("p"), boolV("q"))
			},
			text: "(and p q)",
			kind: KindBoolean,
		},
		{
			name: "Logical disjunction",
			operation: func() *AST {
				return Or(boolV("p"), boolV("q"))
			},
			text: "(or p q)",
			kind: KindBoolean,
		},
		{
			name: "Logical negation",
			operation: func() *AST {
				return Not(boolV("p"))
			},
			text: "(not p)",
			kind: KindBoolean,
		},
		{
			name: "Logical exclusive disjunction",
			operation: func() *AST {
				return Xor(boolV("p"), boolV("q"))
			},
			text: "(xor p q)",
			kind: KindBoolean,
		},
		{
			name: "Logical equivalence",
			operation: func() *AST {
				return Eq(boolV("p"), boolV("q"))
			},
			text: "(= p q)",
			kind: KindBoolean,
		},
		{
			name: "Logical if-then-else",
			operation: func() *AST {
				return ITE(boolV("p"), boolV("q"), boolV("r"))
			},
			text: "(ite p q r)",
			kind: KindBoolean,
		},
		{
			name: "Logical implication",
			operation: func() *AST {
				return Implies(boolV("p"), boolV("q"))
			},
			text: "(=> p q)",
			kind: KindBoolean,
		},
		{
			name: "Equality",
			operation: func() *AST {
				return Eq(boolV("p"), boolV("q"))
			},
			text: "(= p q)",
			kind: KindBoolean,
		},
		{
			name: "Distinction",
			operation: func() *AST {
				return Distinct(boolV("p"), boolV("q"), boolV("r"))
			},
			text: "(distinct p q r)",
			kind: KindBoolean,
		},
		{
			name: "Addition with a real and an integer",
			operation: func() *AST {
				return Add(intV("p"), realV("q"))
			},
			text: "(+ (to_real p) q)",
			kind: KindReal,
		},
		{
			name: "Addition with two reals",
			operation: func() *AST {
				return Add(realV("p"), realV("q"))
			},
			text: "(+ p q)",
			kind: KindReal,
		},
		{
			name: "Addition with two integers",
			operation: func() *AST {
				return Add(intV("p"), intV("q"))
			},
			text: "(+ p q)",
			kind: KindInt,
		},
		{
			name: "Multiplication with a real and an integer",
			operation: func() *AST {
				return Multiply(intV("p"), realV("q"))
			},
			text: "(* (to_real p) q)",
			kind: KindReal,
		},
		{
			name: "Multiplication with two reals",
			operation: func() *AST {
				return Multiply(realV("p"), realV("q"))
			},
			text: "(* p q)",
			kind: KindReal,
		},
		{
			name: "Multiplication with two integers",
			operation: func() *AST {
				return Multiply(intV("p"), intV("q"))
			},
			text: "(* p q)",
			kind: KindInt,
		},
		{
			name: "Subtraction with a real and an integer",
			operation: func() *AST {
				return Subtract(intV("p"), realV("q"))
			},
			text: "(- (to_real p) q)",
			kind: KindReal,
		},
		{
			name: "Subtraction with two reals",
			operation: func() *AST {
				return Subtract(realV("p"), realV("q"))
			},
			text: "(- p q)",
			kind: KindReal,
		},
		{
			name: "Subtraction with two integers",
			operation: func() *AST {
				return Subtract(intV("p"), intV("q"))
			},
			text: "(- p q)",
			kind: KindInt,
		},
		{
			name: "Unary numeric negation (minus) with a real",
			operation: func() *AST {
				return Minus(realV("p"))
			},
			text: "(- p)",
			kind: KindReal,
		},
		{
			name: "Unary numeric negation (minus) with an integer",
			operation: func() *AST {
				return Minus(intV("p"))
			},
			text: "(- p)",
			kind: KindInt,
		},
		{
			name: "Division with a real and an integer",
			operation: func() *AST {
				return Divide(intV("p"), realV("q"))
			},
			text: "(div p (to_int q))",
			kind: KindInt,
		},
		{
			name: "Division with an integer and a real",
			operation: func() *AST {
				return Divide(realV("p"), intV("q"))
			},
			text: "(/ p (to_real q))",
			kind: KindReal,
		},
		{
			name: "Division with two reals",
			operation: func() *AST {
				return Divide(realV("p"), realV("q"))
			},
			text: "(/ p q)",
			kind: KindReal,
		},
		{
			name: "Division with two integers",
			operation: func() *AST {
				return Divide(intV("p"), intV("q"))
			},
			text: "(div p q)",
			kind: KindInt,
		},
		{
			name: "Modulus with a real and an integer",
			operation: func() *AST {
				return Modulus(intV("p"), realV("q"))
			},
			text: "(mod p (to_int q))",
			kind: KindInt,
		},
		{
			name: "Modulus with an integer and a real",
			operation: func() *AST {
				return Modulus(realV("p"), intV("q"))
			},
			text: "(mod (to_int p) q)",
			kind: KindInt,
		},
		{
			name: "Modulus with two reals",
			operation: func() *AST {
				return Modulus(realV("p"), realV("q"))
			},
			text: "(mod (to_int p) (to_int q))",
			kind: KindInt,
		},
		{
			name: "Modulus with two integers",
			operation: func() *AST {
				return Modulus(intV("p"), intV("q"))
			},
			text: "(mod p q)",
			kind: KindInt,
		},
		{
			name: "Remaninder with a real and an integer",
			operation: func() *AST {
				return Remaninder(intV("p"), realV("q"))
			},
			text: "(rem p (to_int q))",
			kind: KindInt,
		},
		{
			name: "Remaninder with an integer and a real",
			operation: func() *AST {
				return Remaninder(realV("p"), intV("q"))
			},
			text: "(rem (to_int p) q)",
			kind: KindInt,
		},
		{
			name: "Remaninder with two reals",
			operation: func() *AST {
				return Remaninder(realV("p"), realV("q"))
			},
			text: "(rem (to_int p) (to_int q))",
			kind: KindInt,
		},
		{
			name: "Remaninder with two integers",
			operation: func() *AST {
				return Remaninder(intV("p"), intV("q"))
			},
			text: "(rem p q)",
			kind: KindInt,
		},
		{
			name: "Power with a real and an integer",
			operation: func() *AST {
				return Power(intV("p"), realV("q"))
			},
			text: "(^ (to_real p) q)",
			kind: KindReal,
		},
		{
			name: "Power with an integer and a real",
			operation: func() *AST {
				return Power(realV("p"), intV("q"))
			},
			text: "(^ p (to_real q))",
			kind: KindReal,
		},
		{
			name: "Power with two reals",
			operation: func() *AST {
				return Power(realV("p"), realV("q"))
			},
			text: "(^ p q)",
			kind: KindReal,
		},
		/*{ This might be machine dependent if it results in an interger or real.
			name: "Power with two integers",
			operation: func() *AST {
				return Power(intV("p"), intV("q"))
			},
			text: "(^ p q)",
			kind: KindInt,
		},*/
		{
			name: "LT with a real and an integer",
			operation: func() *AST {
				return LT(intV("p"), realV("q"))
			},
			text: "(< (to_real p) q)",
			kind: KindBoolean,
		},
		{
			name: "LT with an integer and a real",
			operation: func() *AST {
				return LT(realV("p"), intV("q"))
			},
			text: "(< p (to_real q))",
			kind: KindBoolean,
		},
		{
			name: "LT with two reals",
			operation: func() *AST {
				return LT(realV("p"), realV("q"))
			},
			text: "(< p q)",
			kind: KindBoolean,
		},
		{
			name: "LT with two integers",
			operation: func() *AST {
				return LT(intV("p"), intV("q"))
			},
			text: "(< p q)",
			kind: KindBoolean,
		},
		{
			name: "LE with a real and an integer",
			operation: func() *AST {
				return LE(intV("p"), realV("q"))
			},
			text: "(<= (to_real p) q)",
			kind: KindBoolean,
		},
		{
			name: "LE with an integer and a real",
			operation: func() *AST {
				return LE(realV("p"), intV("q"))
			},
			text: "(<= p (to_real q))",
			kind: KindBoolean,
		},
		{
			name: "LE with two reals",
			operation: func() *AST {
				return LE(realV("p"), realV("q"))
			},
			text: "(<= p q)",
			kind: KindBoolean,
		},
		{
			name: "LE with two integers",
			operation: func() *AST {
				return LE(intV("p"), intV("q"))
			},
			text: "(<= p q)",
			kind: KindBoolean,
		},
		{
			name: "GT with a real and an integer",
			operation: func() *AST {
				return GT(intV("p"), realV("q"))
			},
			text: "(> (to_real p) q)",
			kind: KindBoolean,
		},
		{
			name: "GT with an integer and a real",
			operation: func() *AST {
				return GT(realV("p"), intV("q"))
			},
			text: "(> p (to_real q))",
			kind: KindBoolean,
		},
		{
			name: "GT with two reals",
			operation: func() *AST {
				return GT(realV("p"), realV("q"))
			},
			text: "(> p q)",
			kind: KindBoolean,
		},
		{
			name: "GT with two integers",
			operation: func() *AST {
				return GT(intV("p"), intV("q"))
			},
			text: "(> p q)",
			kind: KindBoolean,
		},
		{
			name: "GE with a real and an integer",
			operation: func() *AST {
				return GE(intV("p"), realV("q"))
			},
			text: "(>= (to_real p) q)",
			kind: KindBoolean,
		},
		{
			name: "GE with an integer and a real",
			operation: func() *AST {
				return GE(realV("p"), intV("q"))
			},
			text: "(>= p (to_real q))",
			kind: KindBoolean,
		},
		{
			name: "GE with two reals",
			operation: func() *AST {
				return GE(realV("p"), realV("q"))
			},
			text: "(>= p q)",
			kind: KindBoolean,
		},
		{
			name: "GE with two integers",
			operation: func() *AST {
				return GE(intV("p"), intV("q"))
			},
			text: "(>= p q)",
			kind: KindBoolean,
		},
		{
			name: "IsInt with an integer",
			operation: func() *AST {
				return IsInt(intV("p"))
			},
			text: "(is_int (to_real p))",
			kind: KindBoolean,
		},
		{
			name: "IsInt with an integer value",
			operation: func() *AST {
				return IsInt(intC(2))
			},
			text: "(is_int (to_real 2))",
			kind: KindBoolean,
		},
	}

	for _, test := range tests {
		// For each test we create a new context such that the factory methods create "fresh" ASTs.
		config := NewConfig()
		context = NewContext(config)
		ast := test.operation()
		kind := ast.Sort().Kind()

		assert.Equal(t, test.text, ast.String(), test.name)
		assert.Equal(t, test.kind, kind, test.name)
	}
}
