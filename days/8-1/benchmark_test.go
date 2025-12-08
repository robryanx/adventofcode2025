package main

import (
	"testing"
)

func BenchmarkSolution(b *testing.B) {
	lines := load()

	for b.Loop() {
		_ = solution(lines)
	}
}
