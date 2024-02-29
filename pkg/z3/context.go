package z3

import "github.com/Brandhoej/gobion/internal/z3"

type Context struct {
	_context *z3.Context
}

func NewContext(config *Config) *Context {
	return &Context{
		_context: z3.NewContext(config._config),
	}
}

func (context *Context) NewSolver() *Solver {
	return &Solver{
		_solver: context._context.NewSolver(),
	}
}

func (context *Context) NewStringSymbol(identifier string) Symbol {
	return Symbol{
		_symbol: context._context.NewStringSymbol(identifier),
	}
}

func (context *Context) NewIntegerSymbol(value int) Symbol {
	return Symbol{
		_symbol: context._context.NewIntegerSymbol(value),
	}
}

type SymbolFactory func(context *Context) Symbol

func WithName(name string) SymbolFactory {
	return func(context *Context) Symbol {
		return context.NewStringSymbol(name)
	}
}

func WithInt(value int) SymbolFactory {
	return func(context *Context) Symbol {
		return context.NewIntegerSymbol(value)
	}
}

func WithSymbol(symbol Symbol) SymbolFactory {
	return func(context *Context) Symbol {
		return symbol
	}
}

func (context *Context) BoolVar(symbolFactory SymbolFactory) Bool {
	sort := context._context.BooleanSort()
	symbol := symbolFactory(context)
	ast := context._context.NewConstant(z3.WithSymbol(symbol._symbol), sort)
	return newBool(ast)
}

func (context *Context) True() Bool {
	return newBool(context._context.NewTrue())
}

func (context *Context) False() Bool {
	return newBool(context._context.NewFalse())
}

func (context *Context) IntVar(symbolFactory SymbolFactory) Int {
	sort := context._context.IntegerSort()
	symbol := symbolFactory(context)
	ast := context._context.NewConstant(z3.WithSymbol(symbol._symbol), sort)
	return newInt(ast)
}

func (context *Context) Int(value int) Int {
	sort := context._context.IntegerSort()
	ast := context._context.NewInt(value, sort)
	return newInt(ast)
}

func (context *Context) RealVar(symbolFactory SymbolFactory) Real {
	sort := context._context.RealSort()
	symbol := symbolFactory(context)
	ast := context._context.NewConstant(z3.WithSymbol(symbol._symbol), sort)
	return newReal(ast)
}

func (context *Context) Real(numerator, denominator int) Real {
	ast := context._context.NewReal(numerator, denominator)
	return newReal(ast)
}
