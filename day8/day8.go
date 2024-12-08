package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"strings"
)

//go:embed input.txt
var input string

func sub(p1 utils.Point, p2 utils.Point) utils.Point {
	return utils.Point{X: p1.X - p2.X, Y: p1.Y - p2.Y}
}

func add(p1 utils.Point, p2 utils.Point) utils.Point {
	return utils.Point{X: p1.X + p2.X, Y: p1.Y + p2.Y}
}

func antinodes(nodes map[byte][]utils.Point, xLen int, yLen int) map[utils.Point]bool {
	ret := make(map[utils.Point]bool)
	for _, points := range nodes {
		for i, point := range points {
			for j, other := range points {
				if i != j {
					distance := sub(point, other)
					antinode := add(point, distance)
					if antinode.X < xLen && antinode.X >= 0 && antinode.Y < yLen && antinode.Y >= 0 {
						ret[antinode] = true
					}
				}
			}
		}
	}
	return ret
}

func antinodesPart2(nodes map[byte][]utils.Point, xLen int, yLen int) map[utils.Point]bool {
	ret := make(map[utils.Point]bool)
	for _, points := range nodes {
		for i, point := range points {
			ret[point] = true // Add all node points
			for j, other := range points {
				if i != j {
					distance := sub(point, other)
					antinode := add(point, distance)
					for {
						if antinode.X < xLen && antinode.X >= 0 && antinode.Y < yLen && antinode.Y >= 0 {
							ret[antinode] = true
							antinode = add(antinode, distance) // Keep adding until we move off the grid
						} else {
							break
						}
					}
				}
			}
		}
	}
	return ret
}

func main() {
	lines := strings.Split(input, "\n")
	//fmt.Println(lines)

	nodes := map[byte][]utils.Point{}
	yLen := len(lines)
	xLen := len(lines[0])
	for y, l := range lines {
		for x, b := range l {
			if b != '.' {
				if nodes[byte(b)] == nil {
					nodes[byte(b)] = []utils.Point{}
				}
				nodes[byte(b)] = append(nodes[byte(b)], utils.Point{X: x, Y: y})
			}
		}
	}
	part1 := antinodes(nodes, xLen, yLen)
	fmt.Println("part1", len(part1))

	part2 := antinodesPart2(nodes, xLen, yLen)
	fmt.Println("part2", len(part2))
}
