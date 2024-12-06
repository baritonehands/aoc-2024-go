package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	lines := strings.Split(input, "\n")
	data := make([][]byte, len(lines))
	start := utils.Point{0, 0}
	for ri, row := range lines {
		data[ri] = []byte(row)
		for ci, c := range row {
			if c == '>' || c == 'v' || c == '<' || c == '^' {
				start.Y = ri
				start.X = ci
			}
		}
	}
	fmt.Println(start, string(data[start.Y][start.X]))

	part1 := func() int {
		cur := start
		seen := map[utils.Point]bool{start: true}
		for {
			//fmt.Println(cur, len(seen))

			switch data[cur.Y][cur.X] {
			case '^':
				if cur.Y == 0 {
					return len(seen)
				}

				if data[cur.Y-1][cur.X] == '#' {
					data[cur.Y][cur.X] = '>'
				} else {
					data[cur.Y][cur.X] = '.'
					cur.Y--
					data[cur.Y][cur.X] = '^'
				}
			case '>':
				if cur.X == len(data[0])-1 {
					return len(seen)
				}

				if data[cur.Y][cur.X+1] == '#' {
					data[cur.Y][cur.X] = 'v'
				} else {

					data[cur.Y][cur.X] = '.'
					cur.X++
					data[cur.Y][cur.X] = '>'
				}
			case 'v':
				if cur.Y == len(data)-1 {
					return len(seen)
				}
				if data[cur.Y+1][cur.X] == '#' {
					data[cur.Y][cur.X] = '<'
				} else {
					data[cur.Y][cur.X] = '.'
					cur.Y++
					data[cur.Y][cur.X] = 'v'
				}
			case '<':
				if cur.X == 0 {
					return len(seen)
				}

				if data[cur.Y][cur.X-1] == '#' {
					data[cur.Y][cur.X] = '^'
				} else {
					data[cur.Y][cur.X] = '.'
					cur.X--
					data[cur.Y][cur.X] = '<'
				}
			}
			seen[cur] = true
		}
		panic("Shouldn't reach here")
	}

	fmt.Println("part1", part1())

	part2 := func(obstacle utils.Point) bool {
		p2data := make([][]byte, len(lines))
		for ri, row := range lines {
			p2data[ri] = []byte(row)
		}
		//fmt.Println(string(p2data[obstacle.X][obstacle.Y]))
		p2data[obstacle.Y][obstacle.X] = '#'

		cur := start
		seen := map[utils.Point]map[byte]bool{start: {p2data[cur.Y][cur.X]: true}}
		for {
			//fmt.Println(cur, len(seen))

			switch p2data[cur.Y][cur.X] {
			case '^':
				if cur.Y == 0 {
					return false
				}

				if p2data[cur.Y-1][cur.X] == '#' {
					p2data[cur.Y][cur.X] = '>'
				} else {
					p2data[cur.Y][cur.X] = '.'
					cur.Y--
					p2data[cur.Y][cur.X] = '^'
				}
			case '>':
				if cur.X == len(p2data[0])-1 {
					return false
				}

				if p2data[cur.Y][cur.X+1] == '#' {
					p2data[cur.Y][cur.X] = 'v'
				} else {

					p2data[cur.Y][cur.X] = '.'
					cur.X++
					p2data[cur.Y][cur.X] = '>'
				}
			case 'v':
				if cur.Y == len(p2data)-1 {
					return false
				}
				if p2data[cur.Y+1][cur.X] == '#' {
					p2data[cur.Y][cur.X] = '<'
				} else {
					p2data[cur.Y][cur.X] = '.'
					cur.Y++
					p2data[cur.Y][cur.X] = 'v'
				}
			case '<':
				if cur.X == 0 {
					return false
				}

				if p2data[cur.Y][cur.X-1] == '#' {
					p2data[cur.Y][cur.X] = '^'
				} else {
					p2data[cur.Y][cur.X] = '.'
					cur.X--
					p2data[cur.Y][cur.X] = '<'
				}
			}
			dirs := seen[cur]
			if dirs == nil {
				seen[cur] = map[byte]bool{p2data[cur.Y][cur.X]: true}
			} else if !dirs[p2data[cur.Y][cur.X]] {
				dirs[p2data[cur.Y][cur.X]] = true
			} else {
				return true // Cycle
			}
		}
		panic("Shouldn't reach here")
	}

	part2Cnt := 0
	for ri, row := range lines {
		for ci := range row {
			if lines[ri][ci] == '.' && (ri != start.Y || ci != start.X) {
				//fmt.Println(utils.Point{X: ci, Y: ri})
				if part2(utils.Point{X: ci, Y: ri}) {
					part2Cnt++
				}
			}
		}
	}
	fmt.Println("part2", part2Cnt)
}
