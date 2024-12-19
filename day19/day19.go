package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

//= `r, wr, b, g, bwu, rb, gb, br
//
//brwrr
//bggr
//gbbr
//rrbgbr
//ubwu
//bwurrg
//brgr
//bbrgwb`

func towelCompare(a, b string) int {
	if len(a) < len(b) {
		return 1
	} else if len(a) > len(b) {
		return -1
	} else {
		return strings.Compare(a, b)
	}
}

type DesignCheck struct {
	designIdx, towelIdx int
}

func nextTowel(towels []string, towelIdx int, designTail string) int {
	for ; towelIdx < len(towels); towelIdx++ {
		towel := towels[towelIdx]
		if strings.HasPrefix(designTail, towel) {
			return towelIdx
		}
	}
	return -1
}

func main() {
	lines := strings.Split(input, "\n")
	towels := strings.Split(lines[0], ", ")
	slices.SortFunc(towels, towelCompare)
	designs := lines[2:]

	fmt.Printf("Max Towel Len: %d\n\n", len(towels[0]))

	part1 := make([]bool, len(designs))
	for i, design := range designs {
		designIdx := 0
		towelIdx := 0
		stack := []DesignCheck{}
		backtracks := 0

		fmt.Printf("Starting: %v\n", design)

		backtrack := func(designTail string) {
			//backtrack
			backtracks++
			lastCheck := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			designIdx = lastCheck.designIdx
			towelIdx = lastCheck.towelIdx + 1
			fmt.Printf("Backtracking: design=%v, (%d,%d), stack=%v\n", designTail, designIdx, towelIdx, stack)
		}

		for {
			if designIdx == len(design) {
				part1[i] = true
				break
			}

			designTail := design[designIdx:]
			foundIdx := nextTowel(towels, towelIdx, designTail)

			if foundIdx == -1 {
				if len(stack) > 0 {
					backtrack(designTail)
				} else {
					// No solution
					break
				}
			} else {
				stack = append(stack, DesignCheck{designIdx, foundIdx})
				designIdx += len(towels[foundIdx])
				towelIdx = 0
			}
		}
		fmt.Printf("Ending: %v, %v, %v\n\n", design, part1[i])
	}

	part1Cnt := it.Len(it.Filter(slices.Values(part1), func(b bool) bool {
		return b
	}))
	fmt.Println("part1", part1Cnt)

	cache := map[string]int{"": 1}
	var dfs func(towels []string, design string) int
	dfs = func(towels []string, design string) int {
		if target, found := cache[design]; found {
			return target
		}

		count := 0
		for _, towel := range towels {
			if strings.HasPrefix(design, towel) {
				count += dfs(towels, design[len(towel):])
			}
		}

		cache[design] = count
		return count
	}

	part2 := it.Fold(slices.Values(designs), func(sum int, design string) int {
		return sum + dfs(towels, design)
	}, 0)
	fmt.Println("part2", part2)
}
