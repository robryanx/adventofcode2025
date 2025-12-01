package main

import (
	"fmt"
	"strconv"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	lines, err := util.ReadStrings("1", false, "\n")
	if err != nil {
		panic(err)
	}

	password := 0
	current := 50
	for line := range lines {
		num, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}

		extraTimes := num / 100

		num %= 100
		prev := current

		switch line[0] {
		case 'L':
			current -= num
			if current < 0 {
				if prev != 0 {
					password++
				}

				current = 100 + current
			}
		case 'R':
			current += num
			if current > 100 && prev != 0 {
				password++
			}
			current %= 100
		}

		if current == 0 {
			password++
		}

		password += extraTimes
	}

	fmt.Println(password)
}
