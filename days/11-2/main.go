package main

import (
	"fmt"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

func main() {
	lines, err := util.ReadStrings("11", false, "\n")
	if err != nil {
		panic(err)
	}

	memo := map[int]int{}
	paths := map[int][][3]byte{}

	for line := range lines {
		parts := strings.Split(line, " ")
		input := parts[0][:len(parts[0])-1]
		inputB := [3]byte{
			input[0],
			input[1],
			input[2],
		}

		for _, path := range parts[1:] {
			key := packKey(inputB, 0)
			paths[key] = append(paths[key], [3]byte{
				path[0],
				path[1],
				path[2],
			})
		}

	}

	fmt.Println(rec(memo, paths, [3]byte{'s', 'v', 'r'}, 0))
}

const (
	hasFft = int8(0b01)
	hasDac = int8(0b10)
)

var (
	fftB = [3]byte{'f', 'f', 't'}
	dacB = [3]byte{'d', 'a', 'c'}
	outB = [3]byte{'o', 'u', 't'}
)

func rec(memo map[int]int, paths map[int][][3]byte, curr [3]byte, flags int8) int {
	matches := 0
	pathKey := packKey(curr, 0)
	for _, path := range paths[pathKey] {
		if compareKey(path, outB) {
			if flags == hasFft|hasDac {
				matches++
			}
		} else {
			nextFlags := flags
			if compareKey(path, fftB) {
				nextFlags |= hasFft
			} else if compareKey(path, dacB) {
				nextFlags |= hasDac
			}

			key := packKey(path, nextFlags)
			if count, ok := memo[key]; ok {
				matches += count
			} else {
				count = rec(memo, paths, path, nextFlags)
				matches += count
				memo[key] = count
			}
		}
	}

	return matches
}

func compareKey(curr, comp [3]byte) bool {
	for i := 0; i < 3; i++ {
		if curr[i] != comp[i] {
			return false
		}
	}

	return true
}

func packKey(path [3]byte, nextFlags int8) int {
	return (int(path[0]-'a') * 10000000) +
		(int(path[1]-'a') * 10000) +
		(int(path[2]-'a') * 10) +
		int(nextFlags)
}
