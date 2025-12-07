package main

import (
	"fmt"
	"iter"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	fmt.Println(solution(load()))
}

func load() iter.Seq[string] {
	lines, err := util.ReadStrings("7", false, "\n")
	if err != nil {
		panic(err)
	}

	return lines
}

func solution(lines iter.Seq[string]) int {
	lineX := map[int]struct{}{}
	splits := 0
	for line := range lines {
		if len(lineX) == 0 {
			for i, ch := range line {
				if ch == 'S' {
					lineX[i] = struct{}{}
					break
				}
			}
		} else {
			for i, ch := range line {
				if ch == '^' {
					if _, ok := lineX[i]; ok {
						lineX[i-1] = struct{}{}
						lineX[i+1] = struct{}{}
						delete(lineX, i)
						splits++
					}
				}
			}
		}
	}

	return splits
}
