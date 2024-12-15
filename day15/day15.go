package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/baritonehands/aoc-2024-go/utils"
	"maps"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

type Input struct {
	grid   [][]byte
	dirs   string
	curPos utils.Point
}

type Input2 Input

func parseInput(lines []string) (Input, Input2) {
	ret := Input{grid: [][]byte{}}
	ret2 := Input2{grid: [][]byte{}}
	dirsAt := 0
	for i, line := range lines {
		if line == "" {
			dirsAt = i
			break
		}
		ret.grid = append(ret.grid, []byte(line))
		ret2.grid = append(ret2.grid, []byte{})
		for _, col := range line {
			if col == 'O' {
				ret2.grid[i] = append(ret2.grid[i], '[', ']')
			} else if col == '@' {
				ret2.grid[i] = append(ret2.grid[i], byte(col), '.')
			} else {
				ret2.grid[i] = append(ret2.grid[i], byte(col), byte(col))
			}
		}

		start := strings.IndexByte(line, '@')
		if start != -1 {
			ret.curPos = utils.Point{Y: i, X: start}
			ret2.curPos = utils.Point{Y: i, X: start * 2}
		}
	}

	ret.dirs = strings.Join(lines[dirsAt:], "")
	ret2.dirs = ret.dirs

	//fmt.Println(ret)
	//fmt.Println(string(ret.grid[ret.curPos.Y][ret.curPos.X]))
	return ret, ret2
}

func (input *Input) move(x, y int) {
	nextY := input.curPos.Y + y
	nextX := input.curPos.X + x

	performMove := func() {
		input.grid[input.curPos.Y][input.curPos.X] = '.'
		input.curPos.X = nextX
		input.curPos.Y = nextY
		input.grid[nextY][nextX] = '@'
	}

	if input.grid[nextY][nextX] == '.' {
		performMove()
	} else if input.grid[nextY][nextX] == 'O' {
		followingY := nextY + y
		followingX := nextX + x
		for i := 0; input.grid[followingY][followingX] != '.' && input.grid[followingY][followingX] != '#'; i++ {
			followingY = followingY + y
			followingX = followingX + x
		}
		if input.grid[followingY][followingX] == '.' {
			performMove()
			input.grid[followingY][followingX] = 'O'
		}
	}
}

type Box struct {
	l, r utils.Point
}

func (input *Input2) move(x, y int) {
	nextY := input.curPos.Y + y
	nextX := input.curPos.X + x

	performMove := func() {
		input.grid[input.curPos.Y][input.curPos.X] = '.'
		input.curPos.X = nextX
		input.curPos.Y = nextY
		input.grid[nextY][nextX] = '@'
	}

	computeVerticalBox := func(cx, cy int) *Box {
		if input.grid[cy+y][cx] == ']' {
			return &Box{utils.Point{Y: cy + y, X: cx - 1}, utils.Point{Y: cy + y, X: cx}}
		} else if input.grid[cy+y][cx] == '[' {
			return &Box{utils.Point{Y: cy + y, X: cx}, utils.Point{Y: cy + y, X: cx + 1}}
		} else {
			return nil
		}
	}

	if input.grid[nextY][nextX] == '.' {
		performMove()
	} else if y == 0 && (input.grid[nextY][nextX] == ']' || input.grid[nextY][nextX] == '[') {
		followingX := nextX + x
		i := 0
		for ; input.grid[nextY][followingX] != '.' && input.grid[nextY][followingX] != '#'; i++ {
			followingX = followingX + x
		}
		if input.grid[nextY][followingX] == '.' {
			for ; i >= 0; i-- {
				input.grid[nextY][followingX] = input.grid[nextY][followingX-x]
				followingX = followingX - x
			}

			performMove()
		}
	} else if x == 0 {
		boxPtr := computeVerticalBox(input.curPos.X, input.curPos.Y)
		if boxPtr != nil {
			box := *boxPtr
			boxes := map[Box]bool{box: true}
			for {
				nextBoxes := it.Fold(maps.Keys(boxes), func(ret map[Box]bool, box Box) map[Box]bool {
					ret[box] = true
					if boxL := computeVerticalBox(box.l.X, box.l.Y); boxL != nil {
						ret[*boxL] = true
					}
					if boxR := computeVerticalBox(box.r.X, box.r.Y); boxR != nil {
						ret[*boxR] = true
					}
					return ret
				}, map[Box]bool{})

				if maps.Equal(boxes, nextBoxes) {
					break
				}
				boxes = nextBoxes
			}

			valid := true
			for toCheck := range maps.Keys(boxes) {
				if input.grid[toCheck.l.Y+y][toCheck.l.X] == '#' || input.grid[toCheck.r.Y+y][toCheck.r.X] == '#' {
					valid = false
					break
				}
			}

			if valid {
				// Ordering is not guaranteed, so just clear all, then set all up/down +/- 1
				for toClear := range maps.Keys(boxes) {
					input.grid[toClear.l.Y][toClear.l.X] = '.'
					input.grid[toClear.r.Y][toClear.r.X] = '.'
				}
				for toMove := range maps.Keys(boxes) {
					input.grid[toMove.l.Y+y][toMove.l.X] = '['
					input.grid[toMove.r.Y+y][toMove.r.X] = ']'
				}
				performMove()
			}
		}

	}
}

func (input Input) gridString() string {
	return strings.Join(slices.Collect(it.Map(slices.Values(input.grid), func(row []byte) string {
		return string(row)
	})), "\n")
}

func main() {
	lines := strings.Split(input, "\n")
	parsed, parsed2 := parseInput(lines)

	fmt.Println(parsed.gridString())
	fmt.Println()

	for i := 0; i < len(parsed.dirs); i++ {
		switch parsed.dirs[i] {
		case '^':
			parsed.move(0, -1)
			break
		case '>':
			parsed.move(1, 0)
			break
		case 'v':
			parsed.move(0, 1)
			break
		case '<':
			parsed.move(-1, 0)
			break
		}
	}

	part1 := 0
	for y, row := range parsed.grid {
		for x := range row {
			if parsed.grid[y][x] == 'O' {
				part1 += 100*y + x
			}

		}
	}

	fmt.Println("part1", part1)

	fmt.Println(Input(parsed2).gridString())
	fmt.Println(parsed2.curPos)
	for i := 0; i < len(parsed2.dirs); i++ {
		switch parsed2.dirs[i] {
		case '^':
			parsed2.move(0, -1)
			break
		case '>':
			parsed2.move(1, 0)
			break
		case 'v':
			parsed2.move(0, 1)
			break
		case '<':
			parsed2.move(-1, 0)
			break
		}
		//fmt.Println(Input(parsed2).gridString())
		//fmt.Println()
	}

	part2 := 0
	for y, row := range parsed2.grid {
		for x := range row {
			if parsed2.grid[y][x] == '[' {
				part2 += 100*y + x
			}
		}
	}

	fmt.Println("part2", part2)
}
