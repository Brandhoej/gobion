package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Equality uint16

const (
	EQ  = Equality(0)
	NEQ = Equality(1)
	LT  = Equality(2)
	LE  = Equality(3)
	GT  = Equality(4)
	GE  = Equality(5)
)

func (equality Equality) Negation() Equality {
	switch equality {
	case EQ:
		return NEQ
	case NEQ:
		return EQ
	case LT:
		return GE
	case LE:
		return GT
	case GT:
		return LE
	case GE:
		return LT
	}
	panic("Unknown equality")
}

func (equality Equality) String() string {
	switch equality {
	case EQ:
		return "="
	case NEQ:
		return "!="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	}
	panic("Unknown equality")
}

func (equality Equality) ast(lhs, rhs *z3.AST) *z3.AST {
	switch equality {
	case EQ:
		return z3.Eq(lhs, rhs)
	case NEQ:
		return z3.Not(z3.Eq(lhs, rhs))
	case LT:
		return z3.LT(lhs, rhs)
	case LE:
		return z3.LE(lhs, rhs)
	case GT:
		return z3.GT(lhs, rhs)
	case GE:
		return z3.GE(lhs, rhs)
	}
	panic("Unknown equality")
}

type Assignment struct {
	symbol     symbols.Symbol
	operator   Equality
	expression *z3.AST
}

func NewAssignment(
	symbol symbols.Symbol, operator Equality, expression *z3.AST,
) Assignment {
	return Assignment{
		symbol:     symbol,
		expression: expression,
		operator:   operator,
	}
}
