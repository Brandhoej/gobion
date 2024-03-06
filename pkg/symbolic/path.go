package symbolic

import "github.com/Brandhoej/gobion/internal/z3"

type Path interface {
	Branch(condition *z3.AST)
}
