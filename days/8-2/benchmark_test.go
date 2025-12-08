package main

import (
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: github.com/robryanx/adventofcode2025/days/8-2
// cpu: AMD Ryzen 7 7800X3D 8-Core Processor
// === RUN   BenchmarkSolution
// BenchmarkSolution
// BenchmarkSolution-16                 132           8123335 ns/op        25233749 B/op       5651 allocs/op
func BenchmarkSolution(b *testing.B) {
	lines := load()

	for b.Loop() {
		_ = solution(lines)
	}
}
