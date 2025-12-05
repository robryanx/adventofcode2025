package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	lines, err := util.ReadStrings("5", false, "\n")
	if err != nil {
		panic(err)
	}

	ranges := [][2]int{}
	total := 0

	isIngredients := false
	for line := range lines {
		if line == "" {
			isIngredients = true
			continue
		}

		if isIngredients {
			id, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}

			for _, r := range ranges {
				if id >= r[0] && id <= r[1] {
					total++
					break
				}
			}
		} else {
			rangeParts := strings.Split(line, "-")
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				panic(err)
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				panic(err)
			}

			ranges = append(ranges, [2]int{start, end})
		}
	}

	fmt.Println(total)
}
