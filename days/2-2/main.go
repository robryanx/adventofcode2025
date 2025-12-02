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
			for j := 2; j <= len(strNum); j++ {
				if len(strNum)%j == 0 {
					matched := true
					curr := strNum[:len(strNum)/j]
					for k := len(strNum) / j; k < len(strNum); k += len(strNum) / j {
						if curr != strNum[k:k+len(strNum)/j] {
							matched = false
							break
						}
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
