package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"math"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//= `029A
//980A
//179A
//456A
//379A`

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
	ret = append(ret, 'A')
	return string(ret)
}

func allPaths(grid Grid, start, end utils.Point) []string {
	if start == end {
		return []string{"A"}
	}

	bfs := [][]utils.Point{{start}}
	ret := []string{}
	shortestLen := math.MaxInt
	for {
		if len(bfs) == 0 {
			break
		}

		curPath := bfs[0]
		bfs = bfs[1:]

		current := curPath[len(curPath)-1]
		for _, neighbor := range current.OrthogonalNeighbors(grid.xMax, grid.yMax) {
			if len(curPath)+1 <= shortestLen && !slices.Contains(curPath, neighbor) && grid.items[neighbor] != 0 {
				nextPath := slices.Clone(curPath)
				nextPath = append(nextPath, neighbor)
				if neighbor == end {
					if len(nextPath) < shortestLen {
						shortestLen = len(nextPath)
						clear(ret)
					}
					ret = append(ret, path2Dirs(nextPath))
				} else {
					bfs = append(bfs, nextPath)
				}
			}
		}
	}
	return ret
}

type Part1 struct {
	code    string
	numeric int64
	presses string
}

func (part1 *Part1) Score() int {
	return int(part1.numeric) * len(part1.presses)
}

func main() {
	codes := strings.Split(input, "\n")
	fmt.Println(codes)

	p1 := allPaths(numberGrid, utils.Point{0, 2}, utils.Point{2, 0})
	p2 := []string{}
	for _, path := range p1[:1] {
		nextBot := allPaths(arrowGrid, arrowGrid.start, arrowGrid.points[path[0]])
		p2 = append(p2, nextBot...)
	}
	fmt.Println(p1, p2)
	fmt.Println()

	fmt.Println(allPaths(numberGrid, utils.Point{0, 0}, utils.Point{0, 1}))
	fmt.Println(allPaths(numberGrid, utils.Point{0, 0}, utils.Point{2, 1}))

	part1 := []Part1{}
	for _, code := range codes {
		sbCode := strings.Builder{}
		lastDigitPos := numberGrid.start
		for _, digit := range code {
			digitPos := numberGrid.points[byte(digit)]
			arrowPaths1 := allPaths(numberGrid, lastDigitPos, digitPos)
			var smallestArrowPath1 *strings.Builder
			for _, arrowPath1 := range arrowPaths1 {
				sbArrow1 := strings.Builder{}
				lastArrow1Pos := arrowGrid.start
				for _, arrow1 := range arrowPath1 {
					arrow1Pos := arrowGrid.points[byte(arrow1)]
					arrowPaths2 := allPaths(arrowGrid, lastArrow1Pos, arrow1Pos)
					var smallestArrowPath2 *strings.Builder
					for _, arrowPath2 := range arrowPaths2 {
						sbArrow2 := strings.Builder{}
						lastArrow2Pos := arrowGrid.start
						for _, arrow2 := range arrowPath2 {
							arrow2Pos := arrowGrid.points[byte(arrow2)]
							arrowPaths3 := allPaths(arrowGrid, lastArrow2Pos, arrow2Pos)
							var smallestArrowPath3 *string
							for _, arrowPath3 := range arrowPaths3 {
								if smallestArrowPath3 == nil || len(arrowPath3) < len(*smallestArrowPath3) {
									smallestArrowPath3 = &arrowPath3
								}
							}
							if smallestArrowPath3 != nil {
								sbArrow2.WriteString(*smallestArrowPath3)
							} else {
								fmt.Println("No smallestArrowPath3", code, string(digit), string(arrow1), string(arrow2))
							}

							lastArrow2Pos = arrow2Pos
						}
						if smallestArrowPath2 == nil || sbArrow2.Len() < smallestArrowPath2.Len() {
							smallestArrowPath2 = &sbArrow2
						}
						lastArrow1Pos = arrow1Pos
					}
					if smallestArrowPath2 != nil {
						sbArrow1.WriteString(smallestArrowPath2.String())
					} else {
						fmt.Println("No smallestArrowPath2", code, string(digit), string(arrow1))
					}
				}
				if smallestArrowPath1 == nil || sbArrow1.Len() < smallestArrowPath1.Len() {
					smallestArrowPath1 = &sbArrow1
				}

				lastDigitPos = digitPos
			}
			if smallestArrowPath1 != nil {
				sbCode.WriteString(smallestArrowPath1.String())
			} else {
				fmt.Println("No smallestArrowPath1", code, string(digit))
			}
		}
		numeric, _ := strconv.ParseInt(code[:len(code)-1], 10, 0)
		result := Part1{code, numeric, sbCode.String()}
		part1 = append(part1, result)
		fmt.Println(len(result.presses), result.numeric, result.Score())
	}
	fmt.Println("part1", part1)
	part1Sum := 0
	for _, item := range part1 {
		part1Sum += item.Score()
	}
	fmt.Println(part1Sum)
}
