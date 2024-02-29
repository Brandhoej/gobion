package z3

import "github.com/Brandhoej/gobion/internal/z3"

type LiftedBoolean z3.LiftedBoolean

func (lbool LiftedBoolean) IsTrue() bool {
	return z3.LiftedBoolean(lbool).IsTrue()
}

func (lbool LiftedBoolean) IsUndefined() bool {
	return z3.LiftedBoolean(lbool).IsUndefined()
}

func (lbool LiftedBoolean) IsFalse() bool {
	return z3.LiftedBoolean(lbool).IsFalse()
}

func (lbool LiftedBoolean) String() string {
	return z3.LiftedBoolean(lbool).String()
}
