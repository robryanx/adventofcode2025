package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/robryanx/adventofcode2025/util"
)

var lineRegex = regexp.MustCompile(`\[(.+)\] (.+) {(.+)}`)

type pos struct {
	currentLights []bool
	depth         int
}

func main() {
	lines, err := util.ReadStrings("10", false, "\n")
	if err != nil {
		panic(err)
	}

	total := 0
	for line := range lines {
		matches := lineRegex.FindStringSubmatch(line)

		targetLights := parseLights(matches[1])
		currentLights := slices.Repeat([]bool{false}, len(targetLights))
		rules := parseRules(strings.Split(matches[2], " "))

		stateLookup := make(map[int]struct{})

		queue := []pos{{
			currentLights: currentLights,
			depth:         0,
		}}

		for len(queue) > 0 {
			next := queue[0]
			queue = queue[1:]

			key := boolsToInt(next.currentLights)
			if _, ok := stateLookup[key]; ok {
				continue
			}
			stateLookup[key] = struct{}{}

			if equal(targetLights, next.currentLights) {
				total += next.depth
				break
			}

			for _, r := range rules {
				lights := slices.Clone(next.currentLights)
				for i := 0; i < len(r); i++ {
					if lights[r[i]] {
						lights[r[i]] = false
					} else {
						lights[r[i]] = true
					}
				}

				queue = append(queue, pos{
					lights,
					next.depth + 1,
				})
			}
		}
	}

	fmt.Println(total)
}

func boolsToInt(b []bool) int {
	var result int
	for i, val := range b {
		if val {
			result |= (1 << i) // Set the i-th bit if val is true
		}
	}
	return result
}

func equal(lightsA, lightsB []bool) bool {
	for i := 0; i < len(lightsA); i++ {
		if lightsA[i] != lightsB[i] {
			return false
		}
	}

	return true
}

func parseLights(lightsStr string) []bool {
	var lights []bool
	for _, ch := range lightsStr {
		if ch == '#' {
			lights = append(lights, true)
		} else {
			lights = append(lights, false)
		}
	}

	return lights
}

func parseRules(rulesStr []string) [][]int {
	var rules [][]int

	for _, rStr := range rulesStr {
		var rule []int
		for lightStr := range strings.SplitSeq(rStr[1:len(rStr)-1], ",") {
			light, err := strconv.Atoi(lightStr)
			if err != nil {
				panic(err)
			}

			rule = append(rule, light)
		}

		rules = append(rules, rule)
	}

	return rules
}
