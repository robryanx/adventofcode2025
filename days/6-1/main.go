package main

import (
	"fmt"
	"iter"
	"strconv"

	"github.com/robryanx/adventofcode2025/util"
)

type operation int32

const (
	operationPlus operation = iota
	operationTimes
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

func solution(lines iter.Seq[string]) int {
	var numberLines [][]int
	var operations []operation

	for line := range lines {
		if line[0] == '*' || line[0] == '+' {
			for _, ch := range line {
				switch ch {
				case '*':
					operations = append(operations, operationTimes)
				case '+':
					operations = append(operations, operationPlus)
				}
			}
		} else {
			var numberLine []int
			var curr []byte
			for i := 0; i < len(line); i++ {
				if line[i] != ' ' {
					curr = append(curr, line[i])
				} else if len(curr) > 0 {
					num, err := strconv.Atoi(string(curr))
					if err != nil {
						panic(err)
					}

					numberLine = append(numberLine, num)
					curr = []byte{}
				}
			}

			if len(curr) > 0 {
				num, err := strconv.Atoi(string(curr))
				if err != nil {
					panic(err)
				}

				numberLine = append(numberLine, num)
			}

			numberLines = append(numberLines, numberLine)
		}
	}

	total := 0
	for i := range len(numberLines[0]) {
		num := numberLines[0][i]
		for j := 1; j < len(numberLines); j++ {
			switch operations[i] {
			case operationPlus:
				num += numberLines[j][i]
			case operationTimes:
				num *= numberLines[j][i]
			}
		}

		total += num
	}

	return total
}
