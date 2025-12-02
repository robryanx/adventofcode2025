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
			strNum := strconv.Itoa(i)
			if len(strNum)%2 == 0 {
				half := len(strNum) / 2
				if strNum[:half] == strNum[half:] {
					total += i
				}
			}
		}
	}

	fmt.Println(total)
}
