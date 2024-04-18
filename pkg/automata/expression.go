package automata

import "github.com/Brandhoej/gobion/internal/z3"

type _Expression struct {
	ast *z3.AST
}

func newExpression(ast *z3.AST) _Expression {
	return _Expression{
		ast: ast,
	}
}

func NewTrue(context *z3.Context) _Expression {
	return newExpression(context.NewTrue())
}

func NewFalse(context *z3.Context) _Expression {
	return newExpression(context.NewFalse())
}

func newEquality(variable *z3.AST, relation Equality, value *z3.AST) _Expression {
	return newExpression(relation.ast(variable, value))
}

func NewConjunction(expression _Expression, expressions ..._Expression) _Expression {
	if len(expressions) == 0 {
		return expression
	}

	return expression.Conjunction(expressions...)
}

func NewDisjunction(expression _Expression, expressions ..._Expression) _Expression {
	if len(expressions) == 0 {
		return expression
	}

	return expression.Disjunction(expressions...)
}

func (expression _Expression) asts(expressions ..._Expression) []*z3.AST {
	asts := make([]*z3.AST, len(expressions)+1)
	asts[0] = expression.ast
	for idx := range expressions {
		asts[idx+1] = expressions[idx].ast
	}
	return asts
}

func (expression _Expression) Conjunction(expressions ..._Expression) _Expression {
	if len(expressions) == 0 {
		return newExpression(expression.ast)
	}

	asts := expression.asts(expressions...)
	return newExpression(z3.And(asts[0], asts[1:]...))
}

func (expression _Expression) Disjunction(expressions ..._Expression) _Expression {
	if len(expressions) == 0 {
		return newExpression(expression.ast)
	}

	asts := expression.asts(expressions...)
	return newExpression(z3.Or(asts[0], asts[1:]...))
}

func (expression _Expression) Negation() _Expression {
	return newExpression(z3.Not(expression.ast))
}

func (expression _Expression) String() string {
	return expression.ast.String()
}
