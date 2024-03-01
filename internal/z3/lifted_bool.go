package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

const (
	LiftedFalse     = LiftedBoolean(C.Z3_L_FALSE)
	LiftedUndefined = LiftedBoolean(C.Z3_L_UNDEF)
	LiftedTrue      = LiftedBoolean(C.Z3_L_TRUE)
)

type LiftedBoolean C.Z3_lbool

func (lbool LiftedBoolean) IsTrue() bool {
	return lbool == LiftedTrue
}

func (lbool LiftedBoolean) IsUndefined() bool {
	return lbool == LiftedUndefined
}

func (lbool LiftedBoolean) IsFalse() bool {
	return lbool == LiftedFalse
}

func (lbool LiftedBoolean) String() string {
	if lbool.IsTrue() {
		return "true"
	} else if lbool.IsFalse() {
		return "false"
	}
	return "undefined"
}
