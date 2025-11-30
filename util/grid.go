package util

import "fmt"

func PrintGrid(grid [][]byte) {
	for y := 0; y < len(grid); y++ {
		fmt.Printf("%s\n", string(grid[y]))
	}
}

func PrintIntGrid(grid [][]int) {
	for y := 0; y < len(grid); y++ {
		fmt.Printf("%d\n", grid[y])
	}
}

func CopyGrid(grid [][]byte, populate bool) [][]byte {
	nextGrid := make([][]byte, len(grid))
	for y := 0; y < len(grid); y++ {
		nextGrid[y] = make([]byte, len(grid[0]))
		if populate {
			for x := 0; x < len(grid[0]); x++ {
				nextGrid[y][x] = grid[y][x]
			}
		}
	}

	return nextGrid
}

func OffsetGridFromGrid(grid [][]byte, offsetY, offsetX int) [][]byte {
	var offsetGrid [][]byte
	for y := offsetY; y < offsetY+8; y++ {
		offsetGrid = append(offsetGrid, grid[y][offsetX:offsetX+8])
	}

	return offsetGrid
}

func AdjacentMatch(grid [][]byte, y, x int, incDiagonal bool, cb func(char byte, y, x int) bool) {
	if y-1 >= 0 {
		earlyExit := cb(grid[y-1][x], y-1, x)
		if earlyExit {
			return
		}

		if incDiagonal {
			if x-1 >= 0 {
				earlyExit := cb(grid[y-1][x-1], y-1, x-1)
				if earlyExit {
					return
				}
			}

			if x+1 < len(grid[0]) {
				earlyExit := cb(grid[y-1][x+1], y-1, x+1)
				if earlyExit {
					return
				}
			}
		}
	}

	if y+1 < len(grid) {
		earlyExit := cb(grid[y+1][x], y+1, x)
		if earlyExit {
			return
		}

		if incDiagonal {
			if x-1 >= 0 {
				earlyExit := cb(grid[y+1][x-1], y+1, x-1)
				if earlyExit {
					return
				}
			}

			if x+1 < len(grid[0]) {
				earlyExit := cb(grid[y+1][x+1], y+1, x+1)
				if earlyExit {
					return
				}
			}
		}
	}

	if x-1 >= 0 {
		earlyExit := cb(grid[y][x-1], y, x-1)
		if earlyExit {
			return
		}
	}

	if x+1 < len(grid[0]) {
		_ = cb(grid[y][x+1], y, x+1)
	}
}
