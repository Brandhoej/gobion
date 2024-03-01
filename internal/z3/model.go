package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"
import "runtime"

type Model struct {
	context *Context
	z3Model C.Z3_model
}

func (context *Context) wrapModel(z3Model C.Z3_model) *Model {
	model := &Model{
		context: context,
		z3Model: z3Model,
	}

	context.do(func() {
		C.Z3_model_inc_ref(context.z3Context, model.z3Model)
	}, z3Model, context)

	runtime.SetFinalizer(model, func(model *Model) {
		context.do(func() {
			C.Z3_model_dec_ref(context.z3Context, model.z3Model)
		}, model)
	})

	return model
}

func (context *Context) NewModel() (model *Model) {
	return context.wrapModel(
		C.Z3_mk_model(context.z3Context),
	)
}

func (solver *Solver) Model() (model *Model) {
	return solver.context.wrapModel(
		C.Z3_solver_get_model(solver.context.z3Context, solver.z3Sovler),
	)
}

// Evaluate the AST node in the given model.
// Return true if succeeded, and the resultant.
//
// If completion is true, then Z3 will assign an interpretation for any constant or function that does
// not have an interpretation in the model. These constants and functions were essentially don't cares.
//
// If completion is false, then Z3 will not assign interpretations to constants for functions that do
// not have interpretations in the model. Evaluation behaves as the identify function in this case.
//
// The evaluation may fail for the following reasons:
// - node contains a quantifier.
// - node is not type correct.
// - the model is partial (that is, the option MODEL_PARTIAL was set to true).
// - Z3_interrupt was invoked during evaluation.
func (model *Model) Eval(node *AST, completion bool) (success bool, resultant *AST) {
	var z3AST C.Z3_ast

	model.context.do(func() {
		success = bool(C.Z3_model_eval(
			model.context.z3Context,
			model.z3Model,
			node.z3AST,
			C.bool(completion),
			&z3AST,
		))
	}, model, node)

	if success {
		resultant = model.context.wrapAST(z3AST)
	}

	return
}

func (model *Model) String() (text string) {
	return compute[string](model.context, func() string {
		return C.GoString(C.Z3_model_to_string(
			model.context.z3Context, model.z3Model,
		))
	}, model)
}
