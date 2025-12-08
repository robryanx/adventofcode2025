package main

import (
	"fmt"
	"iter"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

type point struct {
	x           int
	y           int
	z           int
	circuit     int
	connections []int
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
	var points []point
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

		points = append(points, point{x, y, z, circuit, []int{}})
		circuit++
	}

	connections := 0
	for connections < 1000 {
		midDist := math.MaxInt
		minPosA := -1
		minPosB := -1
		for i := 0; i < len(points); i++ {
			for j := i + 1; j < len(points); j++ {
				checkKey := points[j].x*100000000 + points[j].y*100000 + points[j].z
				if slices.Contains(points[i].connections, checkKey) {
					continue
				}

				dist := (points[i].x-points[j].x)*(points[i].x-points[j].x) +
					(points[i].y-points[j].y)*(points[i].y-points[j].y) +
					(points[i].z-points[j].z)*(points[i].z-points[j].z)

				if dist < midDist {
					midDist = dist
					minPosA = i
					minPosB = j
				}
			}
		}

		if minPosA != -1 {
			aCircuit := points[minPosA].circuit
			bCircuit := points[minPosB].circuit

			for i := 0; i < len(points); i++ {
				if points[i].circuit == bCircuit {
					points[i].circuit = aCircuit
				}
			}

			points[minPosA].connections = append(points[minPosA].connections, points[minPosB].x*100000000+points[minPosB].y*100000+points[minPosB].z)
			points[minPosB].connections = append(points[minPosB].connections, points[minPosA].x*100000000+points[minPosA].y*100000+points[minPosA].z)

			connections++
		}
	}

	sizes := make([]int, len(points))
	for i := 0; i < len(points); i++ {
		count := 0
		for j := 0; j < len(points); j++ {
			if points[j].circuit == i {
				count++
			}
		}

		sizes[i] = count
	}

	slices.Sort(sizes)

	return sizes[len(sizes)-1] * sizes[len(sizes)-2] * sizes[len(sizes)-3]
}
