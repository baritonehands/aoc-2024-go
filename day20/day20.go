package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	pq "github.com/baritonehands/aoc-2024-go/utils/priority_queue"
	"math"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

//= `###############
//#...#...#.....#
//#.#.#.#.#.###.#
//#S#...#.#.#...#
//#######.#.#.###
//#######.#.#...#
//#######.#.###.#
//###..E#...#...#
//###.#######.###
//#...###...#...#
//#.#####.#.###.#
//#.#...#.#.#...#
//#.#.#.#.#.#.###
//#...#...#...###
//###############`

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

type PathCacheKey struct {
	start, end utils.Point
}

var pathCache = map[PathCacheKey][]utils.Point{}

func shortestPath(grid Grid, start, end utils.Point) []utils.Point {
	//if v, ok := pathCache[PathCacheKey{start, end}]; ok {
	//	return v
	//}

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
			ret := walkPath(cameFrom, current)
			//pathCache[PathCacheKey{start, end}] = ret
			return ret
		} else {
			openSet.Poll()

			// For each neighbor of current
			for _, neighbor := range current.OrthogonalNeighbors(grid.xMax, grid.yMax) {
				if !grid.obstacles[neighbor] {
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

func processOneObstacle(ch chan Computation, grid Grid, cheat Cheat, point, start, end utils.Point) {
	path1 := shortestPath(grid, start, point)
	path2 := shortestPath(grid, cheat.end, end)
	diff := len(path1) + len(path2) + nyDistance(cheat.start, cheat.end)
	//if diff > 70 {
	//	fmt.Println(diff)
	//	printBoard(grid, cheat)
	//}
	ch <- Computation{cheat, point, diff}
}

type Computation struct {
	cheat Cheat
	point utils.Point
	diff  int
}

func part2Neighbors(grid Grid, point utils.Point) []utils.Point {
	xStart := point.X - 10
	if xStart < 0 {
		xStart = 0
	}
	yStart := point.Y - 10
	if yStart < 0 {
		yStart = 0
	}

	ret := []utils.Point{}
	for y := yStart; y < point.Y+10 && y <= grid.yMax; y++ {
		for x := xStart; x < point.X+10 && x <= grid.xMax; x++ {
			possible := utils.Point{X: x, Y: y}
			distance := nyDistance(point, possible)
			if distance < 20 && !grid.obstacles[possible] {
				ret = append(ret, possible)
			}
		}
	}
	return ret
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

	baseline := shortestPath(grid, start, end)
	slices.Reverse(baseline)
	part1 := make(chan Computation, 10_000)
	results := 0
	for i, point := range baseline {
		fmt.Println("Iteration", i)
		for _, neighbor1 := range point.OrthogonalNeighbors(grid.xMax, grid.yMax) {
			if grid.obstacles[neighbor1] {
				for _, neighbor2 := range neighbor1.OrthogonalNeighbors(grid.xMax, grid.yMax) {
					if neighbor2 != point && !grid.obstacles[neighbor2] {
						cheat := Cheat{neighbor1, neighbor2}
						go processOneObstacle(part1, grid, cheat, point, start, end)
						results++
					}
				}
			}
		}
	}

	part1Sum := 0
	seen := map[Cheat]bool{}
	part1Freq := map[int]int{}
	for i := 0; i < results; i++ {
		computation := <-part1
		if !seen[computation.cheat] && len(baseline)-computation.diff >= 100 {
			seen[computation.cheat] = true
			part1Freq[len(baseline)-computation.diff]++
			part1Sum++
		}
	}

	fmt.Println("part1Freq", part1Freq)
	fmt.Println("part1", part1Sum)

	part2 := make(chan Computation, 10_000)
	part2Results := 0
	for i, point := range baseline {
		fmt.Println("Iteration", i)
		for _, neighbor1 := range point.OrthogonalNeighbors(grid.xMax, grid.yMax) {
			if grid.obstacles[neighbor1] {
				for _, neighbor2 := range part2Neighbors(grid, neighbor1) {
					if neighbor2 != point {
						cheat := Cheat{neighbor1, neighbor2}
						go processOneObstacle(part2, grid, cheat, point, start, end)
						part2Results++
					}
				}
			}
		}
	}
	part2Sum := 0
	seenPart2 := map[Cheat]bool{}
	part2Freq := map[int]int{}
	for i := 0; i < part2Results; i++ {
		computation := <-part2
		if !seenPart2[computation.cheat] && len(baseline)-computation.diff >= 100 {
			seenPart2[computation.cheat] = true
			part2Freq[len(baseline)-computation.diff]++
			part2Sum++
		}
	}

	fmt.Println("part2Freq", part2Freq)
	fmt.Println("part2", part2Sum)
}
