package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//go:embed input2.txt
var input2 string

func main() {
	input1Lines := strings.Split(input, "\n")
	ordering := make(map[int][]int)
	inverseOrdering := make(map[int][]int)
	for _, line := range input1Lines {
		parts := strings.Split(line, "|")
		from := parts[0]
		to := parts[1]
		fromInt, _ := strconv.Atoi(from)
		toInt, _ := strconv.Atoi(to)

		ordering[fromInt] = append(ordering[fromInt], toInt)
		inverseOrdering[toInt] = append(inverseOrdering[toInt], fromInt)
	}
	fmt.Println(ordering)

	input2Lines := strings.Split(input2, "\n")
	sections := make([][]int, len(input2Lines))
	for i, line := range input2Lines {
		parts := strings.Split(line, ",")
		sections[i] = make([]int, len(parts))

		for j, part := range parts {
			sections[i][j], _ = strconv.Atoi(part)
		}
	}
	fmt.Println(sections)

	isOrdered := func(nums []int) bool {
		last := nums[0]
		for i := 1; i < len(nums); i++ {
			if !slices.Contains(ordering[last], nums[i]) || slices.Contains(inverseOrdering[last], nums[i]) {
				return false
			}
			last = nums[i]
		}
		return true
	}

	sum := 0
	for _, section := range sections {
		if isOrdered(section) {
			sum += section[len(section)/2]
		}
	}
	fmt.Println("part1", sum)

	part2 := 0
	for _, section := range sections {
		if !isOrdered(section) {
			cp := slices.SortedFunc(slices.Values(section), func(lhs int, rhs int) int {
				if slices.Contains(ordering[lhs], rhs) {
					return -1
				} else if slices.Contains(ordering[rhs], lhs) {
					return 1
				} else {
					return 0
				}
			})

			if isOrdered(cp) {
				part2 += cp[len(cp)/2]
			}
		}
	}
	fmt.Println("part2", part2)
}
