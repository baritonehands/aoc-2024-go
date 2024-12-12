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

func parseInput() [][]int {
	lines := strings.Split(input, "\n")
	ret := make([][]int, len(lines))
	for ri, line := range lines {
		row := make([]int, len(line))
		for ci, n := range line {
			row[ci] = int(n - '0')
		}
		ret[ri] = row
	}
	return ret
}

func part1() {
	trails := parseInput()
	xMax := len(trails[0]) - 1
	yMax := len(trails) - 1
	heads := map[utils.Point]bool{}
	tails := map[utils.Point]bool{}
	for ri, row := range trails {
		for ci, col := range row {
			if col == 0 {
				heads[utils.Point{ci, ri}] = true
			} else if col == 9 {
				tails[utils.Point{ci, ri}] = true
			}
		}
	}
	fmt.Println(heads)
	fmt.Println(tails)

	part1Score := 0
	part2Score := 0
	for _, head := range slices.SortedFunc(maps.Keys(heads), utils.PointCompareYX) {
		scored := map[utils.Point]bool{}
		p2Score := 0
		bfs := [][]utils.Point{{head}}
		for len(bfs) > 0 {
			cur := bfs[0]
			bfs = bfs[1:]

			if len(cur) == 10 {
				fmt.Println(head, cur, tails[cur[9]])
				fmt.Println(slices.Collect(it.Map(slices.Values(cur), func(p utils.Point) string {
					return strconv.Itoa(trails[p.Y][p.X])
				})))
				if tails[cur[9]] {
					scored[cur[9]] = true
					p2Score++
				}
			} else {
				end := cur[len(cur)-1]
				for incNeighbor := range it.Filter(slices.Values(end.OrthogonalNeighbors(xMax, yMax)), func(v utils.Point) bool {
					return trails[v.Y][v.X] == trails[end.Y][end.X]+1
				}) {
					next := slices.Clone(cur)
					next = append(next, incNeighbor)
					bfs = append(bfs, next)
				}
			}
		}
		part1Score += len(scored)
		part2Score += p2Score
	}
	fmt.Println("part1", part1Score)
	fmt.Println("part2", part2Score)
}

func main() {
	part1()
}
