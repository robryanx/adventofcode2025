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

		lStart := numLen(start)
		tStart := top(start, lStart-lStart/2)
		lEnd := numLen(start)
		tEnd := top(end, lEnd-lEnd/2)

		for i := tStart; i <= tEnd; i++ {
			j := join(i, i)
			if j >= start && j <= end {
				total += j
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

func join(upper, lower int) int {
	places := numLen(lower)
	for places > 0 {
		upper *= 10
		places--
	}

	return upper + lower
}

func numLen(num int) int {
	count := 0
	for num > 0 {
		count++
		num /= 10
	}

	return count
}
