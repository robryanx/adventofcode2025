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
		best := 0
		var nums []int
		for _, ch := range line {
			nums = append(nums, int(ch-'0'))
		}

		for i := 0; i < len(nums)-1; i++ {
			for j := i + 1; j < len(nums); j++ {
				if nums[i]*10+nums[j] > best {
					best = nums[i]*10 + nums[j]
				}
			}
		}
		total += best
	}

	fmt.Println(total)
}
