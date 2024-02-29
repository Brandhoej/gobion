package z3

import "github.com/Brandhoej/gobion/internal/z3"

type Model struct {
	_model *z3.Model
}

func (solver *Solver) Model() *Model {
	return &Model{
		_model: solver._solver.Model(),
	}
}

func (model *Model) String() string {
	return model._model.String()
}
