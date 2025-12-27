package util

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

type ParamSolution struct {
	NumVars   int
	PivotCols []int
	FreeCols  []int

	Const []*big.Rat   // constant term per variable
	Coeff [][]*big.Rat // coefficients of free vars
}

type bestSolution struct {
	Found bool
	Sum   int64
	B     []int64
	Free  []int64
}

func evaluate(ps ParamSolution, freeVals []int64) ([]int64, bool) {
	if len(freeVals) != len(ps.FreeCols) {
		panic("wrong number of free values")
	}

	b := make([]int64, ps.NumVars)
	for i := 0; i < ps.NumVars; i++ {
		val := new(big.Rat).Set(ps.Const[i])
		for j := range ps.FreeCols {
			term := new(Rat).Mul(ps.Coeff[i][j], NewRat(freeVals[j]))
			val.Add(val, term)
		}
		iv, ok := ToInt64(val)
		if !ok {
			return nil, false
		}
		b[i] = iv
	}
	return b, true
}

func boundT0FromBZero(ineq []inequalityND) (low, high int64, ok bool) {
	low = math.MinInt
	high = math.MaxInt

	for _, in := range ineq {
		if in.Coeff[1].Sign() != 0 {
			continue
		}
		// A*t0 >= R
		if in.Coeff[0].Sign() == 0 {
			// 0 >= R must hold
			if in.R.Sign() > 0 {
				return 0, 0, false
			}
			continue
		}
		div := new(Rat).Quo(in.R, in.Coeff[0]) // R/A (rational)
		if in.Coeff[0].Sign() > 0 {
			// t0 >= ceil(R/A)
			b := CeilRat(div)
			if b > low {
				low = b
			}
		} else {
			// A < 0 => t0 <= floor(R/A)
			b := FloorRat(div)
			if b < high {
				high = b
			}
		}
	}
	return low, high, true
}

func boundT1D(ps ParamSolution, minVal int64) (low, high int64, ok bool) {
	if len(ps.FreeCols) != 1 {
		panic("BoundT1D requires exactly 1 free var")
	}
	low = math.MinInt
	high = math.MaxInt

	minR := NewRat(minVal)

	for i := 0; i < ps.NumVars; i++ {
		a := ps.Coeff[i][0]                  // coefficient of t
		r := new(Rat).Sub(minR, ps.Const[i]) // minVal - const

		if a.Sign() == 0 {
			// const >= minVal must hold
			if r.Sign() > 0 {
				return 0, 0, false
			}
			continue
		}

		div := new(Rat).Quo(r, a) // (minVal-const)/a

		if a.Sign() > 0 {
			b := CeilRat(div)
			if b > low {
				low = b
			}
		} else {
			b := FloorRat(div)
			if b < high {
				high = b
			}
		}
	}
	return low, high, true
}

func boundGen(targetVar int, fixedVals [3]int64, fixedMask [3]bool, ineq []inequalityND) (low, high int64, ok bool) {
	low = math.MinInt64
	high = math.MaxInt64

	for _, in := range ineq {
		// Check if this inequality is usable (only depends on targetVar and fixed vars)
		usable := true
		rhs := new(Rat).Set(in.R)

		for j, c := range in.Coeff {
			if j == targetVar {
				continue
			}
			if c.Sign() == 0 {
				continue
			}

			if fixedMask[j] {
				// Move fixed term to RHS: R - c*val
				term := new(Rat).Mul(c, NewRat(fixedVals[j]))
				rhs.Sub(rhs, term)
			} else {
				// Variable j is not fixed and has non-zero coeff.
				// We cannot use this inequality to bound targetVar yet.
				usable = false
				break
			}
		}

		if !usable {
			continue
		}

		// Solve c*t >= rhs
		c := in.Coeff[targetVar]
		if c.Sign() == 0 {
			if rhs.Sign() > 0 {
				// 0 >= positive -> impossible
				return 0, 0, false
			}
			continue
		}

		div := new(Rat).Quo(rhs, c)
		if c.Sign() > 0 {
			// t >= ceil(div)
			b := CeilRat(div)
			if b > low {
				low = b
			}
		} else {
			// t <= floor(div)
			b := FloorRat(div)
			if b < high {
				high = b
			}
		}
	}
	return low, high, true
}

func eliminateVar(ineqs []inequalityND, varIdx int) []inequalityND {
	var pos, neg, zero []inequalityND
	for _, in := range ineqs {
		s := in.Coeff[varIdx].Sign()
		if s > 0 {
			pos = append(pos, in)
		} else if s < 0 {
			neg = append(neg, in)
		} else {
			zero = append(zero, in)
		}
	}

	next := make([]inequalityND, 0, len(zero)+len(pos)*len(neg))
	next = append(next, zero...)

	for _, p := range pos {
		for _, n := range neg {
			// Combine p (coeff > 0) and n (coeff < 0) to eliminate varIdx
			// p: cP * tk + ... >= RP
			// n: cN * tk + ... >= RN
			// Multipliers: mP = |cN|, mN = cP
			// Result: mP*p + mN*n

			cP := p.Coeff[varIdx]
			cN := n.Coeff[varIdx]
			absCN := new(Rat).Abs(cN)

			newCoeffs := make([]*Rat, len(p.Coeff))
			for i := range newCoeffs {
				term1 := new(Rat).Mul(absCN, p.Coeff[i])
				term2 := new(Rat).Mul(cP, n.Coeff[i])
				newCoeffs[i] = new(Rat).Add(term1, term2)
			}

			term1 := new(Rat).Mul(absCN, p.R)
			term2 := new(Rat).Mul(cP, n.R)
			newR := new(Rat).Add(term1, term2)

			next = append(next, inequalityND{Coeff: newCoeffs, R: newR})
		}
	}
	return next
}

func minimiseSum3D(ps ParamSolution, minVal int64, sumLimit int64, debug bool) bestSolution {
	if len(ps.FreeCols) != 3 {
		panic("MinimiseSum3D requires exactly 3 free vars")
	}

	ineqOriginal := buildIneqND(ps, minVal)
	// We include the sum constraint initially.
	ineqWithSum := addSumConstraint(ineqOriginal, ps, sumLimit)
	// Pointer to the RHS of the sum constraint (last one added)
	sumConstraintRHS := ineqWithSum[len(ineqWithSum)-1].R

	sConst := new(Rat).Set(sumConstraintRHS)
	sConst.Add(sConst, NewRat(sumLimit))

	// 1. Determine order based on independent bounds (Level 1) of the full system
	type varInfo struct {
		id        int
		width     int64
		low       int64
		high      int64
		coeffMass float64
	}
	infos := make([]varInfo, 3)
	emptyVals := [3]int64{}
	emptyMask := [3]bool{}

	// We can use the projection logic to get the EXACT independent range for sorting!
	for i := 0; i < 3; i++ {
		// Eliminate (i+1)%3
		step1 := eliminateVar(ineqWithSum, (i+1)%3)
		// Eliminate (i+2)%3
		step2 := eliminateVar(step1, (i+2)%3)

		// Now step2 has only var i constraints (and constants)
		l, h, ok := boundGen(i, emptyVals, emptyMask, step2)
		width := int64(math.MaxInt64)
		if ok {
			if l < 0 {
				l = 0
			}
			if h > sumLimit {
				h = sumLimit
			}
			if h >= l {
				width = h - l
			}
		} else {
			return bestSolution{} // Impossible
		}

		// Mass heuristic on original inequalities
		mass := 0.0
		for _, in := range ineqWithSum {
			c := in.Coeff[i]
			f, _ := c.Float64()
			mass += math.Abs(f)
		}
		infos[i] = varInfo{id: i, width: width, low: l, high: h, coeffMass: mass}
	}

	// Sort: smallest width first, then largest coeffMass
	p := []int{0, 1, 2}
	less := func(i, j int) bool {
		if infos[i].width != infos[j].width {
			return infos[i].width < infos[j].width
		}
		return infos[i].coeffMass > infos[j].coeffMass
	}
	if less(p[1], p[0]) {
		p[0], p[1] = p[1], p[0]
	}
	if less(p[2], p[1]) {
		p[1], p[2] = p[2], p[1]
	}
	if less(p[1], p[0]) {
		p[0], p[1] = p[1], p[0]
	}

	p0, p1, p2 := p[0], p[1], p[2]

	// 2. Prepare Hierarchical Inequalities
	// Level 3 (Inner): All inequalities (ineqWithSum)
	// Level 2 (Middle): Eliminate p2 from ineqWithSum -> Constraints on p0, p1
	ineqLevel2 := eliminateVar(ineqWithSum, p2)

	// Level 1 (Outer): Eliminate p1 from ineqLevel2 -> Constraints on p0
	ineqLevel1 := eliminateVar(ineqLevel2, p1)

	best := bestSolution{Found: false, Sum: sumLimit + 1}

	var pairsChecked int64
	var t2Checked int64

	start := time.Now()

	// Helper to update bounds based on current best sum
	updateSumConstraint := func(currentBestSum int64) {
		newLimit := currentBestSum
		newLimitRat := NewRat(newLimit)
		newRHS := new(Rat).Sub(sConst, newLimitRat)
		sumConstraintRHS.Set(newRHS)
	}

	// OUTER LOOP (p0) using Level 1 constraints
	t0Low, t0High, ok := boundGen(p0, emptyVals, emptyMask, ineqLevel1)
	if !ok {
		return bestSolution{}
	}
	if t0Low < 0 {
		t0Low = 0
	}
	if t0High > sumLimit {
		t0High = sumLimit
	}

	fixedVals := [3]int64{}
	fixedMask := [3]bool{}

	for val0 := t0Low; val0 <= t0High; val0++ {
		fixedVals[p0] = val0
		fixedMask[p0] = true

		// MIDDLE LOOP (p1) using Level 2 constraints (p2 eliminated)
		t1Low, t1High, ok := boundGen(p1, fixedVals, fixedMask, ineqLevel2)
		if !ok {
			fixedMask[p0] = false
			continue
		}
		if t1Low < 0 {
			t1Low = 0
		}
		if t1High > sumLimit {
			t1High = sumLimit
		}
		if t1Low > t1High {
			fixedMask[p0] = false
			continue
		}

		for val1 := t1Low; val1 <= t1High; val1++ {
			pairsChecked++
			fixedVals[p1] = val1
			fixedMask[p1] = true

			// INNER LOOP (p2) using original constraints
			t2Low, t2High, ok := boundGen(p2, fixedVals, fixedMask, ineqWithSum)
			if !ok {
				fixedMask[p1] = false
				continue
			}
			if t2Low < 0 {
				t2Low = 0
			}
			if t2High > sumLimit {
				t2High = sumLimit
			}
			if t2Low > t2High {
				fixedMask[p1] = false
				continue
			}

			for val2 := t2Low; val2 <= t2High; val2++ {
				t2Checked++
				fixedVals[p2] = val2

				args := []int64{fixedVals[0], fixedVals[1], fixedVals[2]}
				b, ok := evaluate(ps, args)
				if !ok {
					continue
				}

				s, ok := sumAndCheck(b, minVal, sumLimit)
				if !ok {
					continue
				}

				if !best.Found || s < best.Sum {
					best = bestSolution{Found: true, Sum: s, B: b, Free: args}
					updateSumConstraint(s - 1)
				}
			}
			fixedMask[p1] = false
		}
		fixedMask[p0] = false
	}

	if debug {
		fmt.Printf("3D search: order=(t%d,t%d,t%d) (t0,t1) pairs checked=%d, t2 candidates checked=%d, time=%v\n", p0, p1, p2, pairsChecked, t2Checked, time.Since(start))
	}

	return best
}

func addSumConstraint(ineq []inequalityND, ps ParamSolution, sumLimit int64) []inequalityND {
	k := len(ps.FreeCols)
	sConst := new(Rat)
	sCoeff := make([]*Rat, k)
	for j := 0; j < k; j++ {
		sCoeff[j] = new(Rat)
	}

	for i := 0; i < ps.NumVars; i++ {
		sConst.Add(sConst, ps.Const[i])
		for j := 0; j < k; j++ {
			sCoeff[j].Add(sCoeff[j], ps.Coeff[i][j])
		}
	}

	// Constraint: sum(-S_coeff[j]*t[j]) >= S_const - sumLimit
	rhs := new(Rat).Sub(sConst, NewRat(sumLimit))
	finalCoeffs := make([]*Rat, k)
	for j := 0; j < k; j++ {
		finalCoeffs[j] = new(Rat).Neg(sCoeff[j])
	}

	ineq = append(ineq, inequalityND{Coeff: finalCoeffs, R: rhs})
	return ineq
}

func minimiseSum1D(ps ParamSolution, minVal int64, sumLimit int64) bestSolution {
	best := bestSolution{}

	low, high, ok := boundT1D(ps, minVal)
	if !ok {
		return bestSolution{}
	}

	// keep it sane
	if low < 0 {
		low = 0
	}
	if high > sumLimit {
		high = sumLimit
	}
	if low > high {
		return bestSolution{}
	}

	checked := int64(0)

	for t := low; t <= high; t++ {
		checked++
		b, ok := evaluate(ps, []int64{t})
		if !ok {
			continue
		}
		s, ok := sumAndCheck(b, minVal, sumLimit)
		if !ok {
			continue
		}

		if !best.Found || s < best.Sum {
			best = bestSolution{Found: true, Sum: s, B: b, Free: []int64{t}}
		}
	}

	return best
}

type inequalityND struct {
	Coeff []*Rat // length k (free-var coefficients)
	R     *Rat   // RHS
	// meaning: sum_j Coeff[j]*t[j] >= R
}

func buildIneqND(ps ParamSolution, minVal int64) []inequalityND {
	k := len(ps.FreeCols)
	ineq := make([]inequalityND, 0, ps.NumVars)
	minR := NewRat(minVal)

	for i := 0; i < ps.NumVars; i++ {
		r := new(Rat).Sub(minR, ps.Const[i]) // minVal - const[i]
		coeff := make([]*Rat, k)
		for j := 0; j < k; j++ {
			coeff[j] = new(Rat).Set(ps.Coeff[i][j])
		}
		ineq = append(ineq, inequalityND{Coeff: coeff, R: r})
	}
	return ineq
}

func minimiseSum2D(ps ParamSolution, minVal int64, sumLimit int64) bestSolution {
	ineq := buildIneqND(ps, minVal)

	t0Low, t0High, ok := boundT0FromBZero(ineq)
	if !ok {
		return bestSolution{}
	}

	if t0Low < 0 {
		t0Low = 0
	}
	if t0High > sumLimit {
		t0High = sumLimit
	}

	best := bestSolution{Found: false}

	for t0 := t0Low; t0 <= t0High; t0++ {
		// Derive t1 range for this t0
		t1Low := int64(math.MinInt64)
		t1High := int64(math.MaxInt64)

		feasible := true
		t0Rat := NewRat(t0)

		for _, in := range ineq {
			// A*t0 + B*t1 >= R  =>  B*t1 >= R - A*t0
			rhs := new(Rat).Sub(in.R, new(Rat).Mul(in.Coeff[0], t0Rat))

			if in.Coeff[1].Sign() == 0 {
				// already handled in t0 bounds, but safe to check
				if rhs.Sign() > 0 {
					feasible = false
					break
				}
				continue
			}

			div := new(Rat).Quo(rhs, in.Coeff[1]) // (R - A*t0)/B

			if in.Coeff[1].Sign() > 0 {
				// t1 >= ceil(div)
				b := CeilRat(div)
				if b > t1Low {
					t1Low = b
				}
			} else {
				// t1 <= floor(div)
				b := FloorRat(div)
				if b < t1High {
					t1High = b
				}
			}
		}

		if !feasible {
			continue
		}

		// Cap by sumLimit as a safe fallback
		if t1Low < 0 {
			t1Low = 0
		}
		if t1High > sumLimit {
			t1High = sumLimit
		}

		if t1Low > t1High {
			continue
		}

		// Enumerate only feasible t1
		for t1 := t1Low; t1 <= t1High; t1++ {
			b, ok := evaluate(ps, []int64{t0, t1})
			if !ok {
				continue
			}
			s, ok := sumAndCheck(b, minVal, sumLimit)
			if !ok {
				continue
			}
			if !best.Found || s < best.Sum {
				best = bestSolution{Found: true, Sum: s, B: b, Free: []int64{t0, t1}}
			}
		}
	}

	return best
}

func MinimiseSum(ps ParamSolution, minVal int64, sumLimit int64, debug bool) bestSolution {
	switch len(ps.FreeCols) {
	case 0:
		b, ok := evaluate(ps, nil)
		if !ok {
			return bestSolution{}
		}
		s, ok := sumAndCheck(b, minVal, sumLimit)
		if !ok {
			return bestSolution{}
		}
		return bestSolution{Found: true, Sum: s, B: b, Free: nil}
	case 1:
		return minimiseSum1D(ps, minVal, sumLimit)
	case 2:
		return minimiseSum2D(ps, minVal, sumLimit)
	case 3:
		return minimiseSum3D(ps, minVal, sumLimit, debug)
	default:
		panic("too many free vars for this approach (add more pruning or use ILP)")
	}
}

func sumAndCheck(b []int64, minVal int64, sumLimit int64) (int64, bool) {
	sum := int64(0)
	for _, v := range b {
		if v < minVal {
			return 0, false
		}
		sum += v
		if sum > sumLimit {
			return 0, false
		}
	}
	return sum, true
}

func PrintParamSolution(ps ParamSolution) {
	fmt.Println("\nParametric solution:")
	for i := 0; i < ps.NumVars; i++ {
		fmt.Printf("b%d = %s", i, ps.Const[i].RatString())
		for j, fc := range ps.FreeCols {
			c := ps.Coeff[i][j]
			if c.Sign() == 0 {
				continue
			}
			fmt.Printf(" + (%s)*t%d", c.RatString(), fc)
		}
		fmt.Println()
	}
}
