package main

import (
	"fmt"
	"slices"
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
	for line := range lines {
		if line == "" {
			break
		}

		rangeParts := strings.Split(line, "-")
		start, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			panic(err)
		}

		end, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			panic(err)
		}

		found := false
		for i := 0; i < len(ranges); i++ {
			if start <= ranges[i][1] {
				if start > ranges[i][0] {
					start = ranges[i][0]
				}

				pos := i
				for j := i; j < len(ranges); j++ {
					if end >= ranges[j][0] {
						if ranges[j][1] >= end {
							end = ranges[j][1]
						}

						pos = j + 1
					} else {
						break
					}
				}
				next := slices.Clone(ranges[pos:])
				ranges = append(ranges[:i], [2]int{start, end})
				ranges = append(ranges, next...)
				found = true
				break
			}
		}

		if !found {
			ranges = append(ranges, [2]int{start, end})
		}
	}

	total := 0
	for _, r := range ranges {
		total += r[1] - r[0] + 1
	}

	fmt.Println(total)
}
