package main

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

type point struct {
	x       int
	y       int
	z       int
	circuit int
}

type distance struct {
	a        *point
	b        *point
	distance int
}

func main() {
	fmt.Println(solution(load()))
}

func load() iter.Seq[string] {
	lines, err := util.ReadStrings("8", false, "\n")
	if err != nil {
		panic(err)
	}

	return lines
}

func solution(lines iter.Seq[string]) int {
	var points []*point
	circuits := map[int][]*point{}

	circuit := 0
	for line := range lines {
		nums := strings.Split(line, ",")

		x, err := strconv.Atoi(nums[0])
		if err != nil {
			panic(err)
		}

		y, err := strconv.Atoi(nums[1])
		if err != nil {
			panic(err)
		}

		z, err := strconv.Atoi(nums[2])
		if err != nil {
			panic(err)
		}

		p := &point{x, y, z, circuit}
		points = append(points, p)
		circuits[circuit] = append(circuits[circuit], p)
		circuit++
	}

	// precompute distances
	distances := []distance{}
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := (points[i].x-points[j].x)*(points[i].x-points[j].x) +
				(points[i].y-points[j].y)*(points[i].y-points[j].y) +
				(points[i].z-points[j].z)*(points[i].z-points[j].z)

			distances = append(distances, distance{
				a:        points[i],
				b:        points[j],
				distance: dist,
			})
		}
	}

	slices.SortFunc(distances, func(a, b distance) int {
		if a.distance < b.distance {
			return -1
		} else if a.distance > b.distance {
			return 1
		}

		return 0
	})

	for i := 0; i < len(distances); i++ {
		if distances[i].a.circuit == distances[i].b.circuit {
			continue
		}

		aCircuit := distances[i].a.circuit
		bCircuit := distances[i].b.circuit

		for _, p := range circuits[bCircuit] {
			p.circuit = aCircuit
			circuits[aCircuit] = append(circuits[aCircuit], p)
		}
		delete(circuits, bCircuit)

		if len(circuits) == 1 {
			return distances[i].a.x * distances[i].b.x
		}
	}

	return -1
}
