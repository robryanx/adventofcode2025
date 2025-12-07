package main

import (
	"fmt"
	"iter"
	"slices"

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

func solution(linesSeq iter.Seq[string]) int {
	lines := slices.Collect(linesSeq)
	memo := map[int]int{}

	currX := -1
	for i, ch := range lines[0] {
		if ch == 'S' {
			currX = i
			break
		}
	}

	return rec(lines[1:], currX, memo)
}

func rec(remainingLines []string, currX int, memo map[int]int) int {
	counter := 0

	for i, line := range remainingLines {
		if line[currX] == '^' {
			counterAdd := 0
			if c, ok := memo[(currX+1)*100+(len(remainingLines)-i)]; ok {
				counterAdd = c
			} else {
				counterAdd = rec(remainingLines[i+1:], currX+1, memo)
				memo[(currX+1)*100+(len(remainingLines)-i)] = counterAdd
			}
			counter += counterAdd

			currX--
		}
	}

	return counter + 1
}
