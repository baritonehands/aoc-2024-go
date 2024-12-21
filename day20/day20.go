package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	pq "github.com/baritonehands/aoc-2024-go/utils/priority_queue"
	"maps"
	"math"
	"slices"
	"strings"
)

// go:embed input.txt
var input string = `###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`

type Grid struct {
	xMax, yMax int
	obstacles  map[utils.Point]bool
}

func nyDistance(start utils.Point, end utils.Point) int {
	ret := int(math.Abs(float64(end.Y-start.Y)) +
		math.Abs(float64(end.X-start.X)))
	return ret
}

func walkPath(cameFrom map[utils.Point]utils.Point, current utils.Point) []utils.Point {
	ret := []utils.Point{current}
	for {
		if next, ok := cameFrom[current]; ok {
			current = next
			ret = append(ret, current)
		} else {
			break
		}
	}
	return ret
}

func shortestPath(grid Grid, start, end utils.Point, cheat *Cheat) []utils.Point {
	fScore := map[utils.Point]int{start: nyDistance(start, end)}
	fScoreFn := func(point utils.Point) int { return fScore[point] }
	openSet := pq.NewQueue[int, utils.Point](fScoreFn, start)

	cameFrom := map[utils.Point]utils.Point{}
	gScore := map[utils.Point]int{start: 0}

	for {
		if openSet.Len() == 0 {
			return nil
		}

		current := openSet.Peek()

		if current == end {
			// Walk path
			return walkPath(cameFrom, current)
		} else {
			openSet.Poll()

			// For each neighbor of current
			for _, neighbor := range current.OrthogonalNeighbors(grid.xMax, grid.yMax) {
				if (cheat != nil && (cheat.start == neighbor || cheat.end == neighbor)) || !grid.obstacles[neighbor] {
					g := gScore[current] + 1
					gNeighbor, found := gScore[neighbor]
					if !found || g < gNeighbor {
						fScore[neighbor] = g + nyDistance(neighbor, end)

						cameFrom[neighbor] = current
						gScore[neighbor] = g
						openSet.Append(neighbor)
					}
				}
			}
		}

	}

	panic("Shouldn't happen")
}

func nextNeighbor(grid Grid, point utils.Point, dir utils.Point) (bool, utils.Point) {
	neighbor := utils.Point{point.X + dir.X, point.Y + dir.Y}
	if neighbor.X >= 0 && neighbor.X <= grid.xMax && neighbor.Y >= 0 && neighbor.Y <= grid.yMax {
		return true, neighbor
	}
	return false, utils.Point{}
}

type Cheat struct {
	start, end utils.Point
}

func printBoard(grid Grid, cheat Cheat) {
	for y := range grid.yMax + 1 {
		for x := range grid.xMax + 1 {
			point := utils.Point{X: x, Y: y}
			if cheat.start == point {
				fmt.Print("1")
			} else if cheat.end == point {
				fmt.Print("2")
			} else if grid.obstacles[point] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	lines := strings.Split(input, "\n")
	grid := Grid{obstacles: make(map[utils.Point]bool), yMax: len(lines) - 1, xMax: len(lines[0]) - 1}
	start, end := utils.Point{0, 0}, utils.Point{0, 0}
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				grid.obstacles[utils.Point{X: x, Y: y}] = true
			} else if char == 'S' {
				start = utils.Point{X: x, Y: y}
			} else if char == 'E' {
				end = utils.Point{X: x, Y: y}
			}
		}
	}
	fmt.Println(grid)

	baseline := shortestPath(grid, start, end, nil)
	slices.Reverse(baseline)
	part2 := map[Cheat]int{}
	last := baseline[0]
	for _, point := range baseline[1:] {
		for _, neighbor1 := range point.OrthogonalNeighbors(grid.xMax, grid.yMax) {
			found, inPath := nextNeighbor(grid, point, utils.Point{point.X - last.X, point.Y - last.Y})
			if ((found && grid.obstacles[inPath]) || !found) && grid.obstacles[neighbor1] {
				for _, neighbor2 := range neighbor1.OrthogonalNeighbors(grid.xMax, grid.yMax) {
					if neighbor2 != point {
						for _, point2 := range neighbor2.OrthogonalNeighbors(grid.xMax, grid.yMax) {
							if !grid.obstacles[point2] && neighbor1 != point2 && point != point2 {
								cheat := Cheat{neighbor1, neighbor2}
								possible := shortestPath(grid, start, end, &cheat)
								diff := len(baseline) - len(possible)
								if _, found := part2[cheat]; !found && diff > 0 {
									fmt.Println(diff)
									printBoard(grid, cheat)
									part2[cheat] = diff
								}

							}
						}
					}
				}
			}
		}
		last = point
	}

	part2Freq := utils.Frequencies(maps.Values(part2))
	fmt.Println("part2Freq", part2Freq)
	fmt.Println("part2", part2)
}
