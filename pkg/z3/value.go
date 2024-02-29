package z3

import "github.com/Brandhoej/gobion/internal/z3"

type valuation interface {
	sort() *z3.Sort
	ast() *z3.AST
}

func asts[T valuation](values ...T) []*z3.AST {
	asts := make([]*z3.AST, len(values))
	for idx, value := range values {
		asts[idx] = value.ast()
	}
	return asts
}

func NewValuation[T valuation](ast *z3.AST) (out T) {
	var resultant valuation

	switch any(out).(type) {
	case Bool:
		resultant = Bool{
			_ast:  ast,
			_sort: ast.Context().BooleanSort(),
		}
	case Int:
		resultant = Int{
			_ast:  ast,
			_sort: ast.Context().IntegerSort(),
		}
	case Real:
		resultant = Real{
			_ast:  ast,
			_sort: ast.Context().RealSort(),
		}
	}

	return resultant.(T)
}
