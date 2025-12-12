package main

import (
	"fmt"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	lines, err := util.ReadStrings("12-1", true, "\n")
	if err != nil {
		panic(err)
	}

	for line := range lines {
		fmt.Println(line)
	}
}

