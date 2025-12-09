package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	lines, err := util.ReadStrings("9", false, "\n")
	if err != nil {
		panic(err)
	}

	var pairs [][2]int
	for line := range lines {
		valsStr := strings.Split(line, ",")

		valA, err := strconv.Atoi(valsStr[0])
		if err != nil {
			panic(err)
		}
		valB, err := strconv.Atoi(valsStr[1])
		if err != nil {
			panic(err)
		}

		pairs = append(pairs, [2]int{valA, valB})
	}

	largest := 0
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			volume := (abs(pairs[i][0]-pairs[j][0]) + 1) * (abs(pairs[i][1]-pairs[j][1]) + 1)
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
