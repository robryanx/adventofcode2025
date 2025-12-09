package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

var offsets = [][2]int{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

func main() {
	lines, err := util.ReadStrings("9", false, "\n")
	if err != nil {
		panic(err)
	}

	maxY := 0
	maxX := 0

	// arbitrary divisor (happens to be the biggest that works)
	// this makes the problem doable in a relative neive way
	divisor := 750

	var pairs [][4]int
	for line := range lines {
		valsStr := strings.Split(line, ",")

		x, err := strconv.Atoi(valsStr[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(valsStr[1])
		if err != nil {
			panic(err)
		}

		if y/divisor > maxY {
			maxY = y / divisor
		}
		if x/divisor > maxX {
			maxX = x / divisor
		}

		pairs = append(pairs, [4]int{y / divisor, x / divisor, y, x})
	}
	pairs = append(pairs, pairs[0])

	var grid [][]byte
	for y := 0; y < maxY+2; y++ {
		grid = append(grid, bytes.Repeat([]byte{'.'}, maxX+2))
	}

	for i := 0; i < len(pairs)-1; i++ {
		if pairs[i][0] == pairs[i+1][0] {
			for x := min(pairs[i][1], pairs[i+1][1]); x <= max(pairs[i][1], pairs[i+1][1]); x++ {
				grid[pairs[i][0]][x] = 'X'
			}
		} else {
			for y := min(pairs[i][0], pairs[i+1][0]); y < max(pairs[i][0], pairs[i+1][0]); y++ {
				grid[y][pairs[i][1]] = 'X'
			}
		}
	}

	queue := [][2]int{
		{maxY/2 + 10, maxX/2 + 1},
	}

	grid[maxY/2+10][maxX/2+1] = 'X'

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, offset := range offsets {
			y := curr[0]
			x := curr[1]

			if grid[y+offset[0]][x+offset[1]] == '.' {
				grid[y+offset[0]][x+offset[1]] = 'X'
				queue = append(queue, [2]int{y + offset[0], x + offset[1]})
			}
		}
	}

	largest := 0
	for i := 0; i < len(pairs); i++ {
	outer:
		for j := i + 1; j < len(pairs); j++ {
			// this is currently corner cutting and not checking the x-axis
			// of the bounding box
			x := min(pairs[i][1], pairs[j][1])
			for y := min(pairs[i][0], pairs[j][0]); y <= max(pairs[i][0], pairs[j][0]); y++ {
				if grid[y][x] == '.' {
					continue outer
				}
			}

			x = max(pairs[i][1], pairs[j][1])
			for y := min(pairs[i][0], pairs[j][0]); y <= max(pairs[i][0], pairs[j][0]); y++ {
				if grid[y][x] == '.' {
					continue outer
				}
			}

			volume := (abs(pairs[i][2]-pairs[j][2]) + 1) * (abs(pairs[i][3]-pairs[j][3]) + 1)
			if volume > largest {
				largest = volume
			}
		}
	}

	fmt.Println(largest)
}

func abs(val int) int {
	if val < 0 {
		return -val
	}

	return val
}
