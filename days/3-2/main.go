package main

import (
	"fmt"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	lines, err := util.ReadStrings("3", false, "\n")
	if err != nil {
		panic(err)
	}

	total := 0
	for line := range lines {
		var nums []int
		for _, ch := range line {
			nums = append(nums, int(ch-'0'))
		}

		remaining := 12
		pos := -1
		current := 0

		for remaining > 0 {
			best := 0
			bestPos := -1

			for i := pos + 1; i <= len(nums)-remaining; i++ {
				if nums[i] > best {
					best = nums[i]
					bestPos = i
				}
			}

			current += best * pow(remaining-1)
			pos = bestPos

			remaining--
		}

		total += current
	}

	fmt.Println(total)
}

func pow(l int) int {
	c := 1
	for l > 0 {
		c *= 10
		l--
	}

	return c
}
