package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	pairs, err := util.ReadStrings("2", false, ",")
	if err != nil {
		panic(err)
	}

	total := 0
	for pair := range pairs {
		pairParts := strings.Split(pair, "-")

		start, err := strconv.Atoi(pairParts[0])
		if err != nil {
			panic(err)
		}

		end, err := strconv.Atoi(pairParts[1])
		if err != nil {
			panic(err)
		}

		var ranges [][3]int
		curr := start
		for curr < end {
			currL := numLen(curr)
			next := nextLenFirst(currL)

			if next > end {
				ranges = append(ranges, [3]int{currL, curr, end})
			} else {
				ranges = append(ranges, [3]int{currL, curr, next - 1})
			}

			curr = next
		}

		for _, r := range ranges {
			repeat := map[int]struct{}{}
			for i := 2; i <= r[0]; i++ {
				if r[0]%i == 0 {
					checkStart := top(r[1], r[0]-r[0]/i)
					checkEnd := top(r[2], r[0]-r[0]/i)

					for j := checkStart; j <= checkEnd; j++ {
						check := join(j, i)
						if check >= start && check <= end {
							if _, ok := repeat[check]; ok {
								continue
							}

							total += check
							repeat[check] = struct{}{}
						}
					}
				}
			}
		}
	}

	fmt.Println(total)
}

func join(num, times int) int {
	total := 0
	numL := numLen(num)

	for i := numL * times; i > 0; i -= numL {
		tempNum := num
		places := i - numL
		for places > 0 {
			tempNum *= 10
			places--
		}
		total += tempNum
	}

	return total
}

func top(num int, places int) int {
	for places > 0 {
		num /= 10
		places--
	}

	return num
}

func nextLenFirst(l int) int {
	c := 1
	for l > 0 {
		c *= 10
		l--
	}

	return c
}

func numLen(num int) int {
	count := 0
	for num > 0 {
		count++
		num /= 10
	}

	return count
}
