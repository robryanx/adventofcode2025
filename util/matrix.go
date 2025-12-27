package util

import (
	"fmt"
)

type Matrix struct {
	aug [][]*Rat
}

func NewMatrix(aug [][]*Rat) *Matrix {
	return &Matrix{
		aug: aug,
	}
}

// swapRows swaps row i and row j.
func swapRows(aug [][]*Rat, i, j int) {
	aug[i], aug[j] = aug[j], aug[i]
}

// scaleRow multiplies a row by scalar s (s must be non-zero).
func scaleRow(row []*Rat, s *Rat) {
	for k := range row {
		row[k].Mul(row[k], s)
	}
}

// addRowMultiple does: dstRow += factor * srcRow
func addRowMultiple(dstRow, srcRow []*Rat, factor *Rat) {
	for k := range dstRow {
		tmp := new(Rat).Mul(srcRow[k], factor)
		dstRow[k].Add(dstRow[k], tmp)
	}
}

func (m *Matrix) Rref() []int {
	var pivotCols []int

	numRows := len(m.aug)        // number of rows (outputs)
	numVars := len(m.aug[0]) - 1 // number of variables (buttons), last col is RHS

	pivotRow := 0
	for col := 0; col < numVars && pivotRow < numRows; col++ {
		// 1) Find a row >= pivotRow with a non-zero entry in this column.
		sel := -1
		for r := pivotRow; r < numRows; r++ {
			if m.aug[r][col].Sign() != 0 {
				sel = r
				break
			}
		}
		if sel == -1 {
			// No pivot in this column; it becomes a free variable column.
			continue
		}

		// 2) Move selected row to pivotRow position.
		if sel != pivotRow {
			swapRows(m.aug, sel, pivotRow)
		}

		// 3) Scale pivot row so pivot entry becomes 1.
		// pivotVal = aug[pivotRow][col]
		pivotVal := new(Rat).Set(m.aug[pivotRow][col])
		invPivot := new(Rat).Inv(pivotVal) // 1 / pivotVal
		scaleRow(m.aug[pivotRow], invPivot)

		// 4) Eliminate this column in all other rows (make them 0).
		for r := 0; r < numRows; r++ {
			if r == pivotRow {
				continue
			}
			factor := new(Rat).Neg(m.aug[r][col]) // want: row[r] += (-a_rc) * pivotRow
			if factor.Sign() == 0 {
				continue
			}
			addRowMultiple(m.aug[r], m.aug[pivotRow], factor)
		}

		pivotCols = append(pivotCols, col)
		pivotRow++
	}

	return pivotCols
}

func (m *Matrix) ExtractParamSolution(pivotCols []int) ParamSolution {
	numRows := len(m.aug)        // number of rows (outputs)
	numVars := len(m.aug[0]) - 1 // number of variables (buttons), last col is RHS

	freeCols := findFreeCols(numVars, pivotCols)

	// Initialise constants and coefficients
	consts := make([]*Rat, numVars)
	coeffs := make([][]*Rat, numVars)
	for i := 0; i < numVars; i++ {
		consts[i] = NewRat(0)
		coeffs[i] = make([]*Rat, len(freeCols))
		for j := range freeCols {
			coeffs[i][j] = NewRat(0)
		}
	}

	// Free variables: b[free_j] = free_j
	for j, col := range freeCols {
		coeffs[col][j] = NewRat(1)
	}

	// Each pivot row defines one pivot variable
	row := 0
	for _, pc := range pivotCols {
		// Find the row with leading 1 in column pc
		for row < numRows && m.aug[row][pc].Cmp(NewRat(1)) != 0 {
			row++
		}
		if row >= numRows {
			break
		}

		// RHS becomes constant term
		consts[pc] = new(Rat).Set(m.aug[row][numVars])

		// Subtract free variable contributions
		for j, fc := range freeCols {
			// x_pc = rhs - aug[row][fc] * x_fc
			coeffs[pc][j] = new(Rat).Neg(m.aug[row][fc])
		}

		row++
	}

	return ParamSolution{
		NumVars:   numVars,
		PivotCols: pivotCols,
		FreeCols:  freeCols,
		Const:     consts,
		Coeff:     coeffs,
	}
}

func (m *Matrix) IsInconsistent() bool {
	numRows := len(m.aug)        // number of rows (outputs)
	numVars := len(m.aug[0]) - 1 // number of variables (buttons), last col is RHS

	for r := 0; r < numRows; r++ {
		allZero := true
		for c := 0; c < numVars; c++ {
			if m.aug[r][c].Sign() != 0 {
				allZero = false
				break
			}
		}
		if allZero && m.aug[r][numVars].Sign() != 0 {
			return true
		}
	}
	return false
}

func (m *Matrix) Print() {
	for _, row := range m.aug {
		for _, v := range row {
			fmt.Printf("%4s ", v.RatString())
		}
		fmt.Println()
	}
}

func findFreeCols(numVars int, pivotCols []int) []int {
	isPivot := make([]bool, numVars)
	for _, c := range pivotCols {
		isPivot[c] = true
	}

	var free []int
	for c := 0; c < numVars; c++ {
		if !isPivot[c] {
			free = append(free, c)
		}
	}
	return free
}
