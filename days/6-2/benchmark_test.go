package main

import (
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: github.com/robryanx/adventofcode2025/days/6-2
// cpu: AMD Ryzen 7 7800X3D 8-Core Processor
// === RUN   BenchmarkSolution
// BenchmarkSolution
// BenchmarkSolution-16                6464            160997 ns/op          155091 B/op       5785 allocs/op     5800 allocs/op        345 allocs/op
func BenchmarkSolution(b *testing.B) {
	lines := load()

	for b.Loop() {
		_ = solution(lines)
	}
}
