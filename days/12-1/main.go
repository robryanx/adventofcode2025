package main

import (
	"fmt"
	"regexp"

	"github.com/robryanx/adventofcode2025/util"
)

type grid struct {
	north [][]byte
	east  [][]byte
	south [][]byte
	west  [][]byte
}

type requirement struct {
	grid   [][]byte
	shapes []int
}

var requirementRegex = regexp.MustCompile(`([0-9]+)x([0-9]+): ([0-9]+)`)

func main() {
	lines, err := util.ReadStrings("12-1", true, "\n")
	if err != nil {
		panic(err)
	}

	// var grids []grid
	var requirements []requirement

	isGrid := false
	for line := range lines {
		if !isGrid {
			isGrid = true
			continue
		} else if line == "" {
			isGrid = false
			continue
		} else if isGrid {
			matches := requirementRegex.FindStringSubmatch(line)
			if len(matches) == 0 {
				fmt.Println(line)
			}
		}
	}
}
