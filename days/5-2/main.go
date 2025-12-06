package main

import (
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	fmt.Println(solution(load()))
}

func load() iter.Seq[string] {
	lines, err := util.ReadStrings("5", false, "\n")
	if err != nil {
		panic(err)
	}

	return lines
}

func solution(lines iter.Seq[string]) int {
	ranges := make([][2]int, 0, 100)
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

				switch i - pos {
				case -1:
					ranges[i] = [2]int{start, end}
				case 0:
					ranges = append(ranges, [2]int{})
					for j := len(ranges) - 1; j > pos; j-- {
						ranges[j] = ranges[j-1]
					}
					ranges[i] = [2]int{start, end}
				default:
					diff := i - pos
					for j := i + 1; j < len(ranges)+diff+1; j++ {
						ranges[j] = ranges[pos]
						pos++
					}
					ranges[i] = [2]int{start, end}
					ranges = ranges[:len(ranges)+diff+1]
				}

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

	return total
}
