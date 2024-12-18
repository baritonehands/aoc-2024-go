package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/baritonehands/aoc-2024-go/utils"
	pq "github.com/baritonehands/aoc-2024-go/utils/priority_queue"
	"math"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

const xMax = 70
const yMax = 70

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

func safestPath(start, end utils.Point, obstacles map[utils.Point]bool) ([]utils.Point, bool) {
	fScore := map[utils.Point]int{start: nyDistance(start, end)}
	fScoreFn := func(point utils.Point) int { return fScore[point] }
	openSet := pq.NewQueue[int, utils.Point](fScoreFn, start)

	cameFrom := map[utils.Point]utils.Point{}
	gScore := map[utils.Point]int{start: 0}

	for {
		if openSet.Len() == 0 {
			return nil, false
		}

		current := openSet.Peek()

		if current == end {
			// Walk path
			return walkPath(cameFrom, current), true
		} else {
			openSet.Poll()

			// For each neighbor of current
			for _, neighbor := range current.OrthogonalNeighbors(end.X, end.Y) {
				if !obstacles[neighbor] {
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

func printBoard(start, end utils.Point, obstacles map[utils.Point]bool) {
	for y := range end.Y {
		for x := range end.X {
			if obstacles[utils.Point{X: x, Y: y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	lines := strings.Split(input, "\n")
	bytes := []utils.Point{}
	for _, line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		bytes = append(bytes, utils.Point{x, y})
	}
	fmt.Println(bytes)

	part1, _ := safestPath(utils.Point{0, 0}, utils.Point{xMax, yMax}, utils.SeqSet(it.Take(slices.Values(bytes), 1024)))
	fmt.Println("part1", len(part1)-1, part1)

	for i := 1025; i < len(bytes); i++ {
		fmt.Println("Trying", i, bytes[i-1])
		board := utils.SeqSet(it.Take(slices.Values(bytes), uint(i)))
		_, success := safestPath(utils.Point{0, 0}, utils.Point{xMax, yMax}, board)
		if !success {
			printBoard(utils.Point{0, 0}, utils.Point{xMax, yMax}, board)
			break
		}
	}
}
