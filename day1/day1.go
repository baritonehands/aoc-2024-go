package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"math"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	lines := strings.Split(input, "\n")
	left := make([]int, len(lines))
	right := make([]int, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, "   ")
		l, _ := strconv.Atoi(parts[0])
		r, _ := strconv.Atoi(parts[1])
		left[i] = l
		right[i] = r
	}

	slices.Sort(left)
	slices.Sort(right)
	fmt.Printf("%v\n%v\n", left, right)

	sum := 0
	for i, l := range left {
		r := right[i]

		sum += int(math.Abs(float64(l - r)))
	}
	println("part1", sum)

	similarity := 0
	freqs := utils.Frequencies(slices.Values(right))
	for _, l := range left {
		rFreq := freqs[l]
		similarity += l * int(rFreq)
	}

	println("part2", similarity)
}
