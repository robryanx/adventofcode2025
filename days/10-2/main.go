package main

import (
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

var lineRegex = regexp.MustCompile(`\[(.+)\] (.+) {(.+)}`)

type pos struct {
	currentJoltage []int
	depth          int
}

type rule struct {
	presses    []int
	maxPresses int
}

func main() {
	lines, err := util.ReadStrings("10-1", true, "\n")
	if err != nil {
		panic(err)
	}

	total := 0
	for line := range lines {
		matches := lineRegex.FindStringSubmatch(line)

		targetJoltage := parseJoltage(matches[3])
		currentJoltage := slices.Repeat([]int{0}, len(targetJoltage))
		rules := parseRules(targetJoltage, strings.Split(matches[2], " "))

		// printJoltage(targetJoltage)
		// // 0:73 1:56 2:85 3:78 4:79 5:75 6:49 7:84 8:45 9:73
		// currentJoltage = addRule(rules[0].presses, currentJoltage, 26) // [0 1 2 3 4 6 8 9]
		// currentJoltage = addRule(rules[1].presses, currentJoltage, 19) // [1 2 3 4 5 6 9]
		// currentJoltage = addRule(rules[2].presses, currentJoltage, 9)  // [0 1 2 5 7 8 9]
		// currentJoltage = addRule(rules[3].presses, currentJoltage, 1)  // [0 1 4 5 6 7 8]
		// currentJoltage = addRule(rules[4].presses, currentJoltage, 10) // [0 2 3 4 5 7]
		// currentJoltage = addRule(rules[5].presses, currentJoltage, 0)  // [3 4 5 6 7 8]
		// currentJoltage = addRule(rules[6].presses, currentJoltage, 9)  // [2 3 4 7 9]
		// currentJoltage = addRule(rules[7].presses, currentJoltage, 10) // [0 3 5 7 9]
		// currentJoltage = addRule(rules[7].presses, currentJoltage, 0)  // [0 1 4 5 7]
		// currentJoltage = addRule(rules[8].presses, currentJoltage, 0)  // [1 2 3 7 9]
		// currentJoltage = addRule(rules[9].presses, currentJoltage, 0)  // [0 2 3 6 9]
		// currentJoltage = addRule(rules[10].presses, currentJoltage, 0) // [7 8]
		// printJoltage(currentJoltage)
		//
		// for i, j := range currentJoltage {
		// 	fmt.Printf("%d:%d ", i, targetJoltage[i]-j)
		// }
		// fmt.Println("")
		//
		// fmt.Println(rules)

		stateLookup := make(map[string]int)

		presses := rec(rules, 0, 0, currentJoltage, targetJoltage, stateLookup)

		fmt.Println(presses)

		total += presses
	}

	fmt.Println(total)
}

func printJoltage(joltage []int) {
	for i, j := range joltage {
		fmt.Printf("%d:%d ", i, j)
	}
	fmt.Println("")
}

func addRule(rule []int, currentJoltage []int, presses int) []int {
	for _, r := range rule {
		currentJoltage[r] += presses
	}

	return currentJoltage
}

func remaining(currentJoltage, targetJoltage []int) []int {
	var rem []int
	for i, j := range currentJoltage {
		rem = append(rem, targetJoltage[i]-j)
	}
	return rem
}

func rec(rules []rule, presses, depth int, currentJoltage, targetJoltage []int, stateLookup map[string]int) int {
	eq, below, distance := equal(targetJoltage, currentJoltage)
	if eq {
		return presses
	} else if !below {
		return -1
	}

	if len(rules) == 0 {
		return -1
	}
	if !possible(rules, currentJoltage, targetJoltage) {
		return -1
	}

	minJ := 0
	maxJ := maxPresses(rules[0].presses, currentJoltage, targetJoltage)
	if depth == 0 {
		minJ = 17
		maxJ = 23
	}

	if distance < 50 {
		fmt.Println(depth, rules[0], remaining(currentJoltage, targetJoltage), currentJoltage, targetJoltage)
	}

	// key := intsToStr(currentJoltage)
	// min, ok := stateLookup[key]
	// if ok {
	// 	if min < presses {
	// 		return -1
	// 	}
	// }
	// stateLookup[key] = presses

	// fmt.Println(rules[0], currentJoltage, targetJoltage, maxPresses(rules[0].presses, currentJoltage, targetJoltage))
	for i := maxJ; i >= minJ; i-- {
		nextJoltage := slices.Clone(currentJoltage)
		for _, p := range rules[0].presses {
			nextJoltage[p] += i
		}

		presses := rec(rules[1:], presses+i, depth+1, nextJoltage, targetJoltage, stateLookup)
		if presses != -1 {
			return presses
		}
	}

	return -1
}

func possible(rules []rule, currentJoltage, targetJoltage []int) bool {
	types := map[int]struct{}{}
	for _, r := range rules {
		for _, light := range r.presses {
			types[light] = struct{}{}
		}
	}

	for i := 0; i < len(currentJoltage); i++ {
		if targetJoltage[i]-currentJoltage[i] > 0 {
			if _, ok := types[i]; !ok {
				return false
			}
		}
	}

	return true
}

func maxPresses(rule []int, currentJoltage, targetJoltage []int) int {
	minPresses := math.MaxInt
	for i := 0; i < len(rule); i++ {
		if targetJoltage[rule[i]]-currentJoltage[rule[i]] < minPresses {
			minPresses = targetJoltage[rule[i]] - currentJoltage[rule[i]]
		}
	}

	return minPresses
}

func intsToStr(currentJoltage []int) string {
	sb := strings.Builder{}
	sb.Grow(len(currentJoltage))
	for _, j := range currentJoltage {
		sb.WriteString(strconv.Itoa(j))
	}
	return sb.String()
}

func equal(lightsA, lightsB []int) (bool, bool, int) {
	equal := true
	below := true
	distance := 0
	for i := 0; i < len(lightsA); i++ {
		if lightsA[i] != lightsB[i] {
			equal = false
			if lightsA[i] < lightsB[i] {
				below = false
				return false, false, 0
			}
		}
		distance += lightsA[i] - lightsB[i]
	}

	return equal, below, distance
}

func parseJoltage(joltagesStr string) []int {
	var joltages []int
	for joltageStr := range strings.SplitSeq(joltagesStr, ",") {
		joltage, err := strconv.Atoi(joltageStr)
		if err != nil {
			panic(err)
		}

		joltages = append(joltages, joltage)
	}

	return joltages
}

func parseRules(targetJoltage []int, rulesStr []string) []rule {
	var rules []rule
	for _, rStr := range rulesStr {
		minPresses := math.MaxInt
		var r []int
		for lightStr := range strings.SplitSeq(rStr[1:len(rStr)-1], ",") {
			light, err := strconv.Atoi(lightStr)
			if err != nil {
				panic(err)
			}

			if minPresses > targetJoltage[light] {
				minPresses = targetJoltage[light]
			}

			r = append(r, light)
		}

		rules = append(rules, rule{
			presses:    r,
			maxPresses: minPresses,
		})
	}

	slices.SortFunc(rules, func(a, b rule) int {
		if len(a.presses) > len(b.presses) {
			return -1
		} else if len(a.presses) < len(b.presses) {
			return 1
		}

		if a.maxPresses > b.maxPresses {
			return -1
		} else if a.maxPresses < b.maxPresses {
			return 1
		}

		return 0
	})

	return rules
}
