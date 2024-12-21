package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	pq "github.com/baritonehands/aoc-2024-go/utils/priority_queue"
	"math"
	"slices"
	"strconv"
	"strings"
)

// go:embed input.txt
var input string = `029A
980A
179A
456A
379A`

type Grid struct {
	xMax, yMax int
	items      map[utils.Point]byte
	points     map[byte]utils.Point
	start      utils.Point
}

var numberGrid = Grid{
	2, 3,
	map[utils.Point]byte{
		utils.Point{0, 0}: '7',
		utils.Point{1, 0}: '8',
		utils.Point{2, 0}: '9',
		utils.Point{0, 1}: '4',
		utils.Point{1, 1}: '5',
		utils.Point{2, 1}: '6',
		utils.Point{0, 2}: '1',
		utils.Point{1, 2}: '2',
		utils.Point{2, 2}: '3',
		utils.Point{1, 3}: '0',
		utils.Point{2, 3}: 'A',
	},
	map[byte]utils.Point{
		'7': utils.Point{0, 0},
		'8': utils.Point{1, 0},
		'9': utils.Point{2, 0},
		'4': utils.Point{0, 1},
		'5': utils.Point{1, 1},
		'6': utils.Point{2, 1},
		'1': utils.Point{0, 2},
		'2': utils.Point{1, 2},
		'3': utils.Point{2, 2},
		'0': utils.Point{1, 3},
		'A': utils.Point{2, 3},
	},
	utils.Point{2, 3},
}

var arrowGrid = Grid{
	2, 1,
	map[utils.Point]byte{
		utils.Point{1, 0}: '^',
		utils.Point{2, 0}: 'A',
		utils.Point{0, 1}: '<',
		utils.Point{1, 1}: 'v',
		utils.Point{2, 1}: '>',
	},
	map[byte]utils.Point{
		'^': utils.Point{1, 0},
		'A': utils.Point{2, 0},
		'<': utils.Point{0, 1},
		'v': utils.Point{1, 1},
		'>': utils.Point{2, 1},
	},
	utils.Point{2, 0},
}

type NumPath struct {
	from, to int
}

var keypadPaths = map[NumPath]map[byte]int{}

func nyDistance(start utils.Point, end utils.Point) int {
	ret := int(math.Abs(float64(end.Y-start.Y)) +
		math.Abs(float64(end.X-start.X)))
	return ret
}

func path2Dirs(path []utils.Point) string {
	ret := []byte{}
	current := path[0]
	for _, next := range path[1:] {
		var dir byte
		switch next.X - current.X {
		case -1:
			dir = '<'
		case 1:
			dir = '>'
		case 0:
			switch next.Y - current.Y {
			case -1:
				dir = '^'
			case 1:
				dir = 'v'
			}
		}
		ret = append(ret, dir)
		current = next
	}
	return string(ret)
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

func shortestPath(grid Grid, start, end utils.Point) string {
	fScore := map[utils.Point]int{start: nyDistance(start, end)}
	fScoreFn := func(point utils.Point) int { return fScore[point] }
	openSet := pq.NewQueue[int, utils.Point](fScoreFn, start)

	cameFrom := map[utils.Point]utils.Point{}
	gScore := map[utils.Point]int{start: 0}

	for {
		current := openSet.Peek()

		if current == end {
			// Walk path
			path := walkPath(cameFrom, current)
			slices.Reverse(path)
			return path2Dirs(path)
		} else {
			openSet.Poll()

			// For each neighbor of current
			for _, neighbor := range current.OrthogonalNeighbors(grid.xMax, grid.yMax) {
				if grid.items[neighbor] != 0 {
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

type Part1 struct {
	code    string
	numeric int64
	presses string
}

func (part1 *Part1) Score() int {
	return int(part1.numeric) * len(part1.presses)
}

func arrowPath() {

}

func main() {
	codes := strings.Split(input, "\n")
	fmt.Println(codes)

	p1 := shortestPath(numberGrid, utils.Point{0, 2}, utils.Point{2, 0})
	p2 := shortestPath(arrowGrid, arrowGrid.points[p1[0]], arrowGrid.points[p1[1]])
	fmt.Println(p1, p2)
	fmt.Println()

	part1 := []Part1{}
	for _, code := range codes {
		sb := strings.Builder{}
		lastDigitPos := numberGrid.start
		for _, digit := range code {
			digitPos := numberGrid.points[byte(digit)]
			arrowPath1 := shortestPath(numberGrid, lastDigitPos, digitPos)
			arrowPath1 += "A"
			//fmt.Println(code, string(digit), arrowPath1)

			lastArrow1Pos := arrowGrid.start
			for _, arrow1 := range arrowPath1 {
				arrow1Pos := arrowGrid.points[byte(arrow1)]
				arrowPath2 := shortestPath(arrowGrid, lastArrow1Pos, arrow1Pos)
				arrowPath2 += "A"
				//fmt.Println(code, string(digit), string(arrow1), arrowPath2)

				lastArrow2Pos := arrowGrid.start
				for _, arrow2 := range arrowPath2 {
					arrow2Pos := arrowGrid.points[byte(arrow2)]
					arrow3path := shortestPath(arrowGrid, lastArrow2Pos, arrow2Pos)
					arrow3path += "A"
					//fmt.Println(code, string(digit), string(arrow1), string(arrow2), arrow3path)
					sb.WriteString(arrow3path)

					lastArrow2Pos = arrow2Pos
				}

				lastArrow1Pos = arrow1Pos
			}

			lastDigitPos = digitPos
		}
		numeric, _ := strconv.ParseInt(code[:len(code)-1], 10, 0)
		result := Part1{code, numeric, sb.String()}
		part1 = append(part1, result)
		fmt.Println(result.Score())
	}
	fmt.Println("part1", part1)
	part1Sum := 0
	for _, item := range part1 {
		part1Sum += item.Score()
	}
	fmt.Println(part1Sum)
}
