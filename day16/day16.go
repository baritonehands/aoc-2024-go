package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/baritonehands/aoc-2024-go/utils"
	pq "github.com/baritonehands/aoc-2024-go/utils/priority_queue"
	"maps"
	"math"
	"slices"
	"strings"
)

// go:embed input.txt
var input string = `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`

type Grid struct {
	start, end utils.Point
	xMax, yMax int
	obstacles  map[utils.Point]bool
}

type GridMove struct {
	dest   utils.Point
	dir    byte
	weight int
}

func (grid *Grid) Neighbors(p utils.Point, dir byte) []GridMove {
	ret := make([]GridMove, 0, 3)
	if p.X < grid.xMax && !grid.obstacles[utils.Point{p.X + 1, p.Y}] && (dir == 'E' || dir == 'N' || dir == 'S') {
		weight := 1001
		if dir == 'E' {
			weight = 1
		}
		ret = append(ret, GridMove{dest: utils.Point{p.X + 1, p.Y}, dir: 'E', weight: weight})
	}
	if p.Y < grid.yMax && !grid.obstacles[utils.Point{p.X, p.Y + 1}] && (dir == 'W' || dir == 'S' || dir == 'E') {
		weight := 1001
		if dir == 'S' {
			weight = 1
		}
		ret = append(ret, GridMove{dest: utils.Point{p.X, p.Y + 1}, dir: 'S', weight: weight})
	}
	if p.X > 0 && !grid.obstacles[utils.Point{p.X - 1, p.Y}] && (dir == 'W' || dir == 'N' || dir == 'S') {
		weight := 1001
		if dir == 'W' {
			weight = 1
		}
		ret = append(ret, GridMove{dest: utils.Point{p.X - 1, p.Y}, dir: 'W', weight: weight})
	}
	if p.Y > 0 && !grid.obstacles[utils.Point{p.X, p.Y - 1}] && (dir == 'E' || dir == 'N' || dir == 'W') {
		weight := 1001
		if dir == 'N' {
			weight = 1
		}
		ret = append(ret, GridMove{dest: utils.Point{p.X, p.Y - 1}, dir: 'N', weight: weight})
	}
	return ret
}

func nyDistance(start utils.Point, end utils.Point) int {
	ret := int(math.Abs(float64(end.Y-start.Y)) +
		math.Abs(float64(end.X-start.X)))
	return ret
}

func stateString(state [][]int) string {
	sb := strings.Builder{}
	it.ForEach(slices.Values(state), func(row []int) {
		for _, col := range row {
			sb.WriteRune(rune(col + '0'))
		}
		sb.WriteString("\n")
		//sb.WriteString(fmt.Sprintf("%v\n", row))
	})
	return sb.String()
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

func safestPath(grid Grid) ([]utils.Point, int) {
	fScore := map[utils.Point]int{grid.start: nyDistance(grid.start, grid.end)}
	fScoreFn := func(move GridMove) int { return fScore[move.dest] + move.weight }
	openSet := pq.NewQueue[int, GridMove](fScoreFn, GridMove{dest: grid.start, dir: 'E', weight: 1})

	cameFrom := map[utils.Point]utils.Point{}
	gScore := map[utils.Point]int{grid.start: 0}

	for {
		if openSet.Len() == 0 {
			panic("Shouldn't happen")
		}

		current := openSet.Peek()

		if current.dest == grid.end {
			// Walk path
			return walkPath(cameFrom, current.dest), fScore[grid.end]
		} else {
			openSet.Poll()

			// For each neighbor of current
			neighborMoves := grid.Neighbors(current.dest, current.dir)
			for _, neighborMove := range neighborMoves {
				neighbor := neighborMove.dest
				g := gScore[current.dest] + neighborMove.weight
				gNeighbor, found := gScore[neighbor]
				if !found || g < gNeighbor {
					fScore[neighbor] = g + nyDistance(neighbor, grid.end)

					cameFrom[neighbor] = current.dest
					gScore[neighbor] = g
					openSet.Append(neighborMove)
				}
			}
		}

	}

	panic("Shouldn't happen")
}

type Part2Path struct {
	path  []utils.Point
	dir   byte
	score int
}

func matchingPaths(grid Grid, cheapest int) int {
	bfs := []Part2Path{{path: []utils.Point{grid.start}, dir: 'E', score: 0}}
	complete := map[utils.Point]bool{}
	seen := map[utils.Point][]utils.Point{grid.start: {}}
	iterations := 0
	for len(bfs) > 0 {
		current := bfs[0]
		bfs = bfs[1:]
		iterations++

		currentPoint := current.path[len(current.path)-1]
		neighborMoves := grid.Neighbors(currentPoint, current.dir)
		for _, neighborMove := range neighborMoves {
			if seenSlice, found := seen[neighborMove.dest]; found && slices.Contains(seenSlice, currentPoint) {
				continue
			} else if neighborMove.dest == grid.end && current.score+neighborMove.weight == cheapest {
				maps.Insert(complete, maps.All(utils.SeqSet(slices.Values(current.path))))
			} else if current.score+neighborMove.weight <= cheapest {
				nextPath := current.path[0:len(current.path):len(current.path)]
				nextPath = append(nextPath, neighborMove.dest)

				if seen[neighborMove.dest] == nil {
					seen[neighborMove.dest] = []utils.Point{}
				}
				seen[neighborMove.dest] = append(seen[neighborMove.dest], currentPoint)

				bfs = append(bfs, Part2Path{path: nextPath, dir: neighborMove.dir, score: current.score + neighborMove.weight})
			}
		}
	}

	return len(complete) + 1
}

func main() {
	lines := strings.Split(input, "\n")
	grid := Grid{obstacles: make(map[utils.Point]bool), yMax: len(lines) - 1, xMax: len(lines[0]) - 1}
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				grid.obstacles[utils.Point{X: x, Y: y}] = true
			} else if char == 'S' {
				grid.start = utils.Point{X: x, Y: y}
			} else if char == 'E' {
				grid.end = utils.Point{X: x, Y: y}
			}
		}
	}
	fmt.Println(grid)

	fmt.Println(grid.Neighbors(grid.start, 'E'))
	_, part1 := safestPath(grid)
	fmt.Println("part1", part1)

	part2 := matchingPaths(grid, part1)
	fmt.Println("part2", part2)
}
