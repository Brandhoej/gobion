package z3

import "github.com/Brandhoej/gobion/internal/z3"

type Context struct {
	_context *z3.Context
	_true, _false *z3.AST
	_boolean, _integer, _real *z3.Sort
}

func NewContext(config *Config) *Context {
	_context := z3.NewContext(config._config)
	return &Context{
		_context: _context,
		_true: _context.NewBoolean(true),
		_false: _context.NewBoolean(false),
		_boolean: _context.BooleanSort(),
		_integer: _context.IntegerSort(),
		_real: _context.RealSort(),
	}
}

func (context *Context) NewSolver() *Solver {
	return newSolver(context._context.NewSolver())
}

func (context *Context) bool(value bool) (*z3.AST, *z3.Sort) {
	return context._context.NewBoolean(value), context._boolean
}

func (context *Context) int(value int) (*z3.AST, *z3.Sort) {
	return context._context.NewInt(value, context._integer), context._integer
}