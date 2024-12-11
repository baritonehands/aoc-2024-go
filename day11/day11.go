package main

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

var changeCache = map[ChangeStoneCall]int{}

func changeStone(val, n int) int {
	if cached, found := changeCache[ChangeStoneCall{val, n}]; found {
		return cached
	}

	var ret int
	if n == 0 {
		ret = 1
	} else if val == 0 {
		ret = changeStone(1, n-1)
	} else {
		str := strconv.Itoa(val)
		if len(str)%2 == 0 {
			p := int(math.Pow10(len(str) / 2))
			ret = changeStone(val/p, n-1) + changeStone(val%p, n-1)
		} else {
			ret = changeStone(val*2024, n-1)
		}
	}
	changeCache[ChangeStoneCall{val, n}] = ret
	return ret
}

type ChangeStoneCall struct {
	curGen, n int
}

func main() {
	parts := strings.Split(input, " ")
	stones := make([]int, len(parts))
	for i := range stones {
		stones[i], _ = strconv.Atoi(parts[i])
	}
	fmt.Println(stones)

	part1 := 0
	for _, val := range stones {
		part1 += changeStone(val, 25)
	}
	fmt.Println("part1", part1)

	part2 := 0
	for _, val := range stones {
		part2 += changeStone(val, 75)
	}
	fmt.Println("part2", part2)
}
