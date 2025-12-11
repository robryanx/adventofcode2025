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

	paths := map[string][]string{}

	for line := range lines {
		parts := strings.Split(line, " ")
		input := parts[0][:len(parts[0])-1]

		paths[input] = append(paths[input], parts[1:]...)
	}

	fmt.Println(rec(paths, "you"))
}

func rec(paths map[string][]string, curr string) int {
	matches := 0
	for _, path := range paths[curr] {
		if path == "out" {
			matches++
		} else {
			matches += rec(paths, path)
		}
	}

	return matches
}
