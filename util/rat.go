package util

import "math/big"

type Rat = big.Rat

func NewRat(v int64) *Rat {
	return new(Rat).SetInt64(v)
}

func FloorRat(x *Rat) int64 {
	// floor(num/den)
	q := new(big.Int)
	r := new(big.Int)
	q.QuoRem(x.Num(), x.Denom(), r) // truncated toward zero

	// If x is negative and remainder != 0, subtract 1 to get floor.
	if x.Sign() < 0 && r.Sign() != 0 {
		q.Sub(q, big.NewInt(1))
	}
	return q.Int64()
}

func CeilRat(x *Rat) int64 {
	// ceil(x) = -floor(-x)
	neg := new(Rat).Neg(x)
	return -FloorRat(neg)
}

func ToInt64(r *Rat) (int64, bool) {
	if !isWhole(r) {
		return 0, false
	}
	return r.Num().Int64(), true
}

func isWhole(r *Rat) bool {
	return r.Denom().Cmp(big.NewInt(1)) == 0
}
