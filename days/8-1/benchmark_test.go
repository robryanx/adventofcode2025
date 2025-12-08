package main

import (
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: github.com/robryanx/adventofcode2025/days/8-1
// cpu: AMD Ryzen 7 7800X3D 8-Core Processor
// === RUN   BenchmarkSolution
// BenchmarkSolution
// BenchmarkSolution-16                 157           7079943 ns/op        24358797 B/op       4972 allocs/op
func BenchmarkSolution(b *testing.B) {
	lines := load()

	for b.Loop() {
		_ = solution(lines)
	}
}
