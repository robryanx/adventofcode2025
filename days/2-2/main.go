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

		for i := start; i <= end; i++ {
			l := numLen(i)
			for j := 2; j <= l; j++ {
				if l%j == 0 {
					matched := true
					match := top(i, l-l/j)
					if match == bottom(i, l/j) {
						if j > 2 {
							for k := l / j; k < l-l/j; k += l / j {
								next := bottom(top(i, (l-k-l/j)), l/j)
								if next != match {
									matched = false
									break
								}
							}
						}
					} else {
						matched = false
					}

					if matched {
						total += i
						break
					}
				}
			}
		}
	}

	fmt.Println(total)
}

func top(num int, places int) int {
	for places > 0 {
		num /= 10
		places--
	}

	return num
}

func bottom(num int, places int) int {
	c := 1
	for places > 0 {
		c *= 10
		places--
	}

	return num % c
}

func numLen(num int) int {
	count := 0
	for num > 0 {
		count++
		num /= 10
	}

	return count
}
