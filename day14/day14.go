package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/baritonehands/aoc-2024-go/utils"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

const xLen = 101
const yLen = 103

type Robot struct {
	pos      utils.Point
	velocity utils.Point
}

func (robot *Robot) move() {
	robot.pos.X += robot.velocity.X
	robot.pos.Y += robot.velocity.Y
	if robot.pos.X >= xLen {
		robot.pos.X -= xLen
	} else if robot.pos.X < 0 {
		robot.pos.X += xLen
	}
	if robot.pos.Y >= yLen {
		robot.pos.Y -= yLen
	} else if robot.pos.Y < 0 {
		robot.pos.Y += yLen
	}
}

func parsePoint(s string) utils.Point {
	parts := strings.Split(s[2:], ",")
	ret := utils.Point{}
	ret.X, _ = strconv.Atoi(parts[0])
	ret.Y, _ = strconv.Atoi(parts[1])
	return ret
}

func printRobots(robots []Robot) {
	freq := utils.Frequencies(it.Map(slices.Values(robots), func(r Robot) utils.Point {
		return r.pos
	}))
	for y := range yLen {
		for x := range xLen {
			if n, found := freq[utils.Point{X: x, Y: y}]; found {
				fmt.Print(n)
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
	robots := make([]Robot, len(lines))

	for i, line := range lines {
		robot := Robot{}
		parts := strings.Split(line, " ")
		robot.pos = parsePoint(parts[0])
		robot.velocity = parsePoint(parts[1])
		robots[i] = robot
		fmt.Println(robot)
	}
	//printRobots(robots)

	qWidth := xLen / 2
	qHeight := yLen / 2
	iterateSeconds := func(breakIf func(i int, central int, pointCnt int) bool) int {
		points := make(map[utils.Point]bool, len(robots))
		for i := 0; ; i++ {
			central := 0
			clear(points)
			for ir := range robots {
				robots[ir].move()
				points[robots[ir].pos] = true
				if robots[ir].pos.X == qWidth || robots[ir].pos.Y == qHeight {
					central++
				}
			}
			if breakIf(i, central, len(points)) {
				return i
			}
		}
	}
	iterateSeconds(func(i int, central int, maxFreq int) bool {
		return i == 100
	})

	quad1 := []utils.Point{}
	quad2 := []utils.Point{}
	quad3 := []utils.Point{}
	quad4 := []utils.Point{}

	for _, robot := range robots {
		if robot.pos.X < qWidth && robot.pos.Y < qHeight {
			quad1 = append(quad1, robot.pos)
		} else if robot.pos.X > qWidth && robot.pos.Y < qHeight {
			quad2 = append(quad2, robot.pos)
		} else if robot.pos.X > qWidth && robot.pos.Y > qHeight {
			quad3 = append(quad3, robot.pos)
		} else if robot.pos.X < qWidth && robot.pos.Y > qHeight {
			quad4 = append(quad4, robot.pos)
		}
	}
	fmt.Println("part1", len(quad1)*len(quad2)*len(quad3)*len(quad4))

	limit := 100_000_000 - 100
	iterations := 100
	fmt.Println(iterateSeconds(func(i int, central int, pointCnt int) bool {
		iterations++
		if i%100_000 == 0 {
			fmt.Println(i)
		}
		return i >= limit || pointCnt == len(robots)
	}))
	printRobots(robots)
	fmt.Println(iterations + 1)
}
