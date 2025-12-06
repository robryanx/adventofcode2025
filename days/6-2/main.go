package main

import (
	"fmt"
	"iter"
	"slices"
	"strconv"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	fmt.Println(solution(load()))
}

func load() iter.Seq[string] {
	lines, err := util.ReadStrings("6", false, "\n")
	if err != nil {
		panic(err)
	}

	return lines
}

func solution(linesSeq iter.Seq[string]) int {
	total := 0
	lines := slices.Collect(linesSeq)

	last := len(lines[0]) + 1
	for i := len(lines[0]) - 1; i >= 0; i-- {
		if lines[len(lines)-1][i] == '*' || lines[len(lines)-1][i] == '+' {
			// extract lines
			var sums []string
			for j := 0; j < len(lines)-1; j++ {
				sums = append(sums, lines[j][i:last-1])
			}

			lineTotal := 0
			for k := 0; k < len(sums[0]); k++ {
				var buildNum []byte
				for x := 0; x < len(sums); x++ {
					if sums[x][k] != ' ' {
						buildNum = append(buildNum, sums[x][k])
					}
				}

				num, err := strconv.Atoi(string(buildNum))
				if err != nil {
					panic(err)
				}

				switch lines[len(lines)-1][i] {
				case '+':
					lineTotal += num
				case '*':
					if lineTotal == 0 {
						lineTotal = num
					} else {
						lineTotal *= num
					}
				}
			}

			total += lineTotal
			last = i
		}
	}

	return total
}
