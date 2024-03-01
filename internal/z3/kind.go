package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

// Kind is a Z3 type.
type Kind int

// The different kinds of Z3 types.
const (
	KindUninterpreted     = Kind(C.Z3_UNINTERPRETED_SORT)
	KindBoolean           = Kind(C.Z3_BOOL_SORT)
	KindInt               = Kind(C.Z3_INT_SORT)
	KindReal              = Kind(C.Z3_REAL_SORT)
	KindBitVector         = Kind(C.Z3_BV_SORT)
	KindArray             = Kind(C.Z3_ARRAY_SORT)
	KindDatatype          = Kind(C.Z3_DATATYPE_SORT)
	KindRelation          = Kind(C.Z3_RELATION_SORT)
	KindFiniteDomain      = Kind(C.Z3_FINITE_DOMAIN_SORT)
	KindFloatingPoint     = Kind(C.Z3_FLOATING_POINT_SORT)
	KindRoundingMode      = Kind(C.Z3_ROUNDING_MODE_SORT)
	KindSequence          = Kind(C.Z3_SEQ_SORT)
	KindRegularExpression = Kind(C.Z3_RE_SORT)
	KindCharacter         = Kind(C.Z3_CHAR_SORT)
	KindTypeVariable      = Kind(C.Z3_TYPE_VAR)
	KindUnknown           = Kind(C.Z3_UNKNOWN_SORT)
)
