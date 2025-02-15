package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type ClawConfig struct {
	A, B  utils.Point
	Prize utils.Point
}

func decodePair(pair string) utils.Point {
	parts := strings.Split(pair, ", ")
	X, _ := strconv.Atoi(parts[0][2:])
	Y, _ := strconv.Atoi(parts[1][2:])
	return utils.Point{X: X, Y: Y}
}

func isDivisible(config ClawConfig, factor int) (bool, int) {
	for lhs := config.A.X; lhs < config.Prize.X+factor; lhs += config.A.X {
		rhs := config.Prize.X + factor - lhs
		if rhs%config.B.X == 0 {
			div1 := lhs / config.A.X
			div2 := rhs / config.B.X

			if factor == 0 && (div1 > 100 || div2 > 100) {
				continue
			}

			if div1*config.A.Y+div2*config.B.Y == config.Prize.Y+factor {
				return true, 3*div1 + 1*div2
			}
		}
	}
	return false, 0
}

func part2Solve(aX, aY, bX, bY, pX, pY, costA, costB int) int {
	bNumer := pX*aY - pY*aX
	bDenom := aY*bX - aX*bY
	if bNumer%bDenom != 0 {
		return 0
	}
	b := bNumer / bDenom
	aNumer := pX - (b * bX)
	aDenom := aX

	if aNumer%aDenom != 0 {
		return 0
	}
	a := aNumer / aDenom

	return a*costA + b*costB
}

func main() {
	lines := strings.Split(input, "\n")
	fmt.Println(lines)

	part1 := 0
	part2 := 0
	groups := utils.Partition(lines, 3, 4)
	configs := make([]ClawConfig, 0)
	for parts := range groups {
		config := ClawConfig{}
		a := parts[0]
		config.A = decodePair(a[10:])
		b := parts[1]
		config.B = decodePair(b[10:])
		prize := parts[2]
		config.Prize = decodePair(prize[7:])
		configs = append(configs, config)

		divP1, tokensP1 := isDivisible(config, 0)
		if divP1 {
			part1 += tokensP1
		}

		sol1 := part2Solve(config.A.X, config.A.Y, config.B.X, config.B.Y, config.Prize.X+10000000000000, config.Prize.Y+10000000000000, 3, 1)
		sol2 := part2Solve(config.B.X, config.B.Y, config.A.X, config.A.Y, config.Prize.X+10000000000000, config.Prize.Y+10000000000000, 1, 3)
		if sol1 > 0 && sol2 > 0 {
			if sol1 < sol2 {
				part2 += sol1
			} else {
				part2 += sol2
			}
		} else {
			if sol1 > sol2 {
				part2 += sol1
			} else {
				part2 += sol2
			}
		}
		//fmt.Printf("%v, %v, %v -> %v, %v\n", config.A, config.B, config.Prize, divP1, tokensP1)
	}
	fmt.Println("part1", part1)
	fmt.Println("part2", part2)
}
