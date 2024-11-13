package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"maps"
)

//go:embed input.txt
var input string

func main() {
	ret := utils.PartitionFunc2([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, func(n int) int {
		return n / 3
	})

	for v := range ret {
		fmt.Printf("%v\n", maps.Collect(v))
	}
}
