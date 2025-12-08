package main

import (
	"fmt"
	"iter"
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
	distances := make([]distance, 0, len(points)*len(points))
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			distX := points[i].x - points[j].x
			distY := points[i].y - points[j].y
			distZ := points[i].z - points[j].z

			distances = append(distances, distance{
				a:        points[i],
				b:        points[j],
				distance: (distX * distX) + (distY * distY) + (distZ * distZ),
			})
		}
	}

	pq := util.NewPriorityQueue(distances, func(a, b distance) bool {
		return a.distance < b.distance
	})

	for {
		dist, ok := pq.Pop()
		if !ok {
			return -1
		}
		if dist.a.circuit == dist.b.circuit {
			continue
		}

		bCircuit := dist.b.circuit

		for _, p := range circuits[bCircuit] {
			p.circuit = dist.a.circuit
			circuits[dist.a.circuit] = append(circuits[dist.a.circuit], p)
		}
		delete(circuits, bCircuit)

		if len(circuits) == 1 {
			return dist.a.x * dist.b.x
		}
	}
}
