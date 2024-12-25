package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

//= `#####
//.####
//.####
//.####
//.#.#.
//.#...
//.....
//
//#####
//##.##
//.#.##
//...##
//...#.
//...#.
//.....
//
//.....
//#....
//#....
//#...#
//#.#.#
//#.###
//#####
//
//.....
//.....
//#.#..
//###..
//###.#
//###.#
//#####
//
//.....
//.....
//.....
//#....
//#.#..
//#.#.#
//#####`

func main() {
	lines := strings.Split(input, "\n")
	lines = append(lines, "\n")
	keys := [][]int{}
	locks := [][]int{}
	for keySeq := range utils.Partition(lines, 8, 8) {
		key := slices.Collect(keySeq)
		isLock := key[0][0] == '#'
		sizes := make([]int, len(key[0]))
		for _, row := range key[1 : len(key)-2] {
			for col, val := range row {
				if val == '#' {
					sizes[col]++
				}
			}
		}
		if isLock {
			locks = append(locks, sizes)
		} else {
			keys = append(keys, sizes)
		}
	}

	part1 := 0
	for _, lock := range locks {
		for _, key := range keys {
			success := true
			for i := 0; i < len(lock); i++ {
				if key[i]+lock[i] > 5 {
					success = false
					break
				}
			}
			if success {
				part1++
			}
		}
	}
	fmt.Println("part1", part1)
}
