package main

import (
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: github.com/robryanx/adventofcode2025/days/7-2
// cpu: AMD Ryzen 7 7800X3D 8-Core Processor
// === RUN   BenchmarkSolution
// BenchmarkSolution
// BenchmarkSolution-16                4212            257094 ns/op          102746 B/op        176 allocs/op
func BenchmarkSolution(b *testing.B) {
	lines := load()

	for b.Loop() {
		_ = solution(lines)
	}
}
