package main

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

var lineRegex = regexp.MustCompile(`\[(.+)\] (.+) {(.+)}`)
var debug = false

func main() {
	lines, err := util.ReadStrings("10", false, "\n")
	if err != nil {
		panic(err)
	}

	total := 0
	for line := range lines {
		matches := lineRegex.FindStringSubmatch(line)

		targetJoltage, targetSum := parseJoltage(matches[3])
		rules := parseRules(strings.Split(matches[2], " "))

		m := util.NewMatrix(buildAugmentedMatrix(rules, targetJoltage))

		pivots := m.Rref()
		if m.IsInconsistent() {
			panic("Inconsistent")
		}

		ps := m.ExtractParamSolution(pivots)

		best := util.MinimiseSum(ps, 0, int64(targetSum), debug)
		if !best.Found {
			panic("No feasible solution under constraints")
		}

		total += int(best.Sum)

		if debug {
			m.Print()
			fmt.Println(pivots)
			fmt.Println("\nBest solution:")
			fmt.Println("Free vars:", best.Free, "for columns", ps.FreeCols)
			fmt.Println("Sum:", best.Sum)
			for i, v := range best.B {
				fmt.Printf("b%d = %d\n", i, v)
			}
		}
	}

	fmt.Println(total)
}

func buildAugmentedMatrix(rules [][]int, target []int) [][]*big.Rat {
	numOutputs := len(target)
	numButtons := len(rules)

	// Allocate matrix: outputs × (buttons + RHS)
	aug := make([][]*big.Rat, numOutputs)
	for i := 0; i < numOutputs; i++ {
		aug[i] = make([]*big.Rat, numButtons+1)
		for j := 0; j < numButtons+1; j++ {
			aug[i][j] = util.NewRat(0)
		}
	}

	// Fill A matrix from rules
	for i, rule := range rules {
		for _, button := range rule {
			aug[button][i] = util.NewRat(1)
		}
	}

	// Fill RHS column
	for i, val := range target {
		aug[i][numButtons] = util.NewRat(int64(val))
	}

	return aug
}

func parseJoltage(joltagesStr string) ([]int, int) {
	var joltages []int
	sum := 0
	for joltageStr := range strings.SplitSeq(joltagesStr, ",") {
		joltage, err := strconv.Atoi(joltageStr)
		if err != nil {
			panic(err)
		}

		joltages = append(joltages, joltage)
		sum += joltage
	}

	return joltages, sum
}

func parseRules(rulesStr []string) [][]int {
	var rules [][]int
	for _, rStr := range rulesStr {
		var r []int
		for lightStr := range strings.SplitSeq(rStr[1:len(rStr)-1], ",") {
			light, err := strconv.Atoi(lightStr)
			if err != nil {
				panic(err)
			}

			r = append(r, light)
		}

		rules = append(rules, r)
	}

	return rules
}
