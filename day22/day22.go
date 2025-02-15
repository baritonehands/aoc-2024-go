package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/baritonehands/aoc-2024-go/utils"
	"maps"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//= `1
//2
//3
//2024`

//= `1
//10
//100
//2024`

func mix(n int, value int) int {
	return n ^ value
}

func prune(n int) int {
	return n % 16777216
}

func nextSecretNumber(n int) int {
	op1 := prune(mix(n, 64*n))
	op2 := prune(mix(op1, op1/32))
	return prune(mix(op2, op2*2048))
}

type Diff struct {
	diff, price int8
}

type WindowKey struct {
	n1, n2, n3, n4 int8
}

type Window struct {
	key   WindowKey
	price int8
}

func main() {
	lines := strings.Split(input, "\n")
	numbers := make([]int, len(lines))
	for i, line := range lines {
		numbers[i], _ = strconv.Atoi(line)
	}

	for i := 0; i < 2000; i++ {
		for j, n := range numbers {
			numbers[j] = nextSecretNumber(n)
		}
	}
	part1Sum := 0
	for _, n := range numbers {
		part1Sum += n
	}
	fmt.Println("part1", part1Sum)

	allNumbers := make([][]int, len(lines))
	for i, line := range lines {
		allNumbers[i] = make([]int, 2001)
		allNumbers[i][0], _ = strconv.Atoi(line)
		for j := 1; j < len(allNumbers[i]); j++ {
			allNumbers[i][j] = nextSecretNumber(allNumbers[i][j-1])
		}
	}

	allDiffs := make([][]Diff, len(allNumbers))
	for i, nums := range allNumbers {
		allDiffs[i] = make([]Diff, len(nums))
		j := 0
		for l, r := range it.Zip(slices.Values(nums), slices.Values(nums[1:])) {
			allDiffs[i][j] = Diff{int8((r % 10) - (l % 10)), int8(r % 10)}
			j++
		}
	}

	allWindows := make([][]Window, 0)
	for _, diffs := range allDiffs {
		windows := slices.Collect(it.Map(utils.Partition(diffs, 4, 1), func(slice []Diff) Window {
			if len(slice) == 4 {
				return Window{WindowKey{slice[0].diff, slice[1].diff, slice[2].diff, slice[3].diff}, slice[3].price}
			} else {
				return Window{}
			}
		}))
		allWindows = append(allWindows, windows)
	}

	bestBuys := map[WindowKey]map[int]int8{}
	for i, windows := range allWindows {
		for _, window := range windows {
			if bestBuys[window.key] == nil {
				bestBuys[window.key] = map[int]int8{}
			}
			if _, found := bestBuys[window.key][i]; !found {
				bestBuys[window.key][i] = window.price
			}
		}
	}
	//fmt.Println(bestBuys)

	best := 0
	for _, freqs := range bestBuys {
		//fmt.Printf("%v\n", freqs)
		sum := 0
		for v := range maps.Values(freqs) {
			sum += int(v)
		}
		if sum > best {
			best = sum
		}
	}

	fmt.Println("part2", best)
}
