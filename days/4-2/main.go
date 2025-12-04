package main

import (
	"fmt"

	"github.com/robryanx/adventofcode2025/util"
)

var offsets = [][2]int{
	{-1, 0},
	{-1, -1},
	{-1, 1},
	{1, 0},
	{1, 1},
	{1, -1},
	{0, -1},
	{0, 1},
}

func main() {
	rows, err := util.ReadStrings("4", false, "\n")
	if err != nil {
		panic(err)
	}

	var grid [][]byte
	for row := range rows {
		grid = append(grid, []byte(row))
	}

	total := 0

	for {
		removable := 0
		for y := range len(grid) {
			for x := range len(grid[0]) {
				if grid[y][x] != '@' {
					continue
				}

				count := 0
				movable := true
				for _, offset := range offsets {
					if y+offset[0] < 0 ||
						y+offset[0] > len(grid)-1 ||
						x+offset[1] < 0 ||
						x+offset[1] > len(grid[0])-1 {
						continue
					}
					if grid[y+offset[0]][x+offset[1]] == '@' {
						count++
						if count > 3 {
							movable = false
							break
						}
					}
				}

				if movable {
					grid[y][x] = '.'
					removable++
				}
			}
		}

		if removable == 0 {
			break
		}
		total += removable
	}

	fmt.Println(total)
}
