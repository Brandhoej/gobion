package symbolic

import "github.com/Brandhoej/gobion/internal/z3"

type Valuations interface {
	Load(symbol Symbol) (*z3.AST, bool)
	Store(symbol Symbol, value *z3.AST)
}

type ValuationsMap struct {
	valuations map[Symbol]*z3.AST
}

func NewEnvironmentMap() *ValuationsMap {
	return &ValuationsMap{
		valuations: map[Symbol]*z3.AST{},
	}
}

func (environment *ValuationsMap) Load(symbol Symbol) (value *z3.AST, exists bool) {
	value, exists = environment.valuations[symbol]
	return
}

func (environment *ValuationsMap) Store(symbol Symbol, value *z3.AST) {
	environment.valuations[symbol] = value
}
