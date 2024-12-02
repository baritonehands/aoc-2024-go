package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

// This didn't end up working
func increasing(levels []int, max int, badAllowed int) (bool, []int) {
	last := levels[0]
	bad := 0
	badIdx := -1
	start := 1
	if levels[1] <= levels[0] {
		last = levels[1]
		bad = levels[0]
		badIdx = 0
		start = 2
	}

	if bad > badAllowed {
		return false, []int{}
	}

	for i := start; i < len(levels); i++ {
		if levels[i] > last && levels[i]-last <= max {
			last = levels[i]
		} else if bad < badAllowed {
			bad = levels[i]
			badIdx = i
		} else {
			return false, []int{}
		}
	}
	return true, []int{badIdx, bad}
}

func decreasing(levels []int, max int, badAllowed int) (bool, []int) {
	cp := make([]int, len(levels))
	copy(cp, levels)
	slices.Reverse(cp)
	return increasing(cp, max, badAllowed)
}

func isMonotone(levels []int) bool {
	cp := make([]int, len(levels))
	copy(cp, levels)
	slices.Sort(cp)
	if slices.Equal(cp, levels) {
		return true
	}

	slices.Reverse(cp)
	return slices.Equal(cp, levels)
}

func isGradual(levels []int) bool {
	lhs := levels[:len(levels)-1]
	rhs := levels[1:]
	fmt.Printf("%v <-> %v\n", lhs, rhs)
	for l, r := range it.Zip(slices.Values(lhs), slices.Values(rhs)) {
		if (l > r && l-r > 3) || (r > l && r-l > 3) {
			fmt.Println("false")
			return false
		} else if l == r {
			return false
		}
	}
	fmt.Println("true")
	return true
}

func main() {
	lines := strings.Split(input, "\n")
	levels := make([][]int, len(lines))
	part1 := 0
	for i, line := range lines {
		parts := strings.Split(line, " ")
		levels[i] = make([]int, len(parts))
		for j, part := range parts {
			levels[i][j], _ = strconv.Atoi(part)
		}

		if isGradual(levels[i]) && isMonotone(levels[i]) {
			part1++
		}
	}
	println("part1", part1)

	part2 := 0
	for _, row := range levels {
		for i := range len(row) {
			test := slices.Concat(row[:i], row[i+1:])
			if isGradual(test) && isMonotone(test) {
				part2++
				fmt.Printf("%v -> %v\n", row, test)
				break
			}
		}
	}
	println("part2", part2)
}
