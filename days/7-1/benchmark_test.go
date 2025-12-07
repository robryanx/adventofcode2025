package main

import (
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: github.com/robryanx/adventofcode2025/days/7-1
// cpu: AMD Ryzen 7 7800X3D 8-Core Processor
// === RUN   BenchmarkSolution
// BenchmarkSolution
// BenchmarkSolution-16               10000            102818 ns/op           29406 B/op        158 allocs/op
func BenchmarkSolution(b *testing.B) {
	lines := load()

	for b.Loop() {
		_ = solution(lines)
	}
}
