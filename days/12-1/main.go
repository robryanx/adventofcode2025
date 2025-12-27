package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

type requirement struct {
	width  int
	height int
	shapes []int
}

var requirementRegex = regexp.MustCompile(`([0-9]+)x([0-9]+): ([0-9\s]+)`)

func main() {
	lines, err := util.ReadStrings("12", false, "\n")
	if err != nil {
		panic(err)
	}

	var grids [][][]byte
	var requirements []requirement
	var currGrid [][]byte

	isGrid := false
	for line := range lines {
		matches := requirementRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			var shapes []int
			for shapeStr := range strings.SplitSeq(matches[3], " ") {
				shape, err := strconv.Atoi(shapeStr)
				if err != nil {
					panic(err)
				}

				shapes = append(shapes, shape)
			}

			width, err := strconv.Atoi(matches[1])
			if err != nil {
				panic(err)
			}

			height, err := strconv.Atoi(matches[2])
			if err != nil {
				panic(err)
			}

			requirements = append(requirements, requirement{
				width, height, shapes,
			})
		} else if !isGrid {
			isGrid = true
			continue
		} else if line == "" {
			grids = append(grids, currGrid)
			currGrid = nil
			isGrid = false
			continue
		} else if isGrid {
			currGrid = append(currGrid, []byte(line))
		}
	}

	var gridCounts []int
	for _, grid := range grids {
		count := 0
		for y := 0; y < len(grid); y++ {
			for x := 0; x < len(grid[0]); x++ {
				if grid[y][x] == '#' {
					count++
				}
			}
		}
		gridCounts = append(gridCounts, count)
	}

	// interestingly just using the min area of the presents works
	// prsumably this would depend on input
	total := 0
	for _, req := range requirements {
		area := req.height * req.width

		shapeTotal := 0
		for i := 0; i < len(req.shapes); i++ {
			shapeTotal += gridCounts[i] * req.shapes[i]
		}
		if shapeTotal <= area {
			total++
		}
	}

	fmt.Println(total)
}
