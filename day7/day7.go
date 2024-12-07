package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type Part1 struct {
	result   int
	operands []int
}

var operatorCache = make(map[int][][]byte)

func operators(n int, zeroValue []byte) [][]byte {
	if operatorCache[n] != nil {
		return operatorCache[n]
	} else if n == 1 {
		ret := make([][]byte, 0, len(zeroValue))
		for _, v := range zeroValue {
			ret = append(ret, []byte{v})
		}
		operatorCache[n] = ret
		return ret
	} else {
		prev := operators(n-1, zeroValue)
		ret := make([][]byte, 0, len(prev)*2)
		for _, ops := range prev {
			for _, op := range zeroValue {
				toAdd := slices.Clone(ops)
				toAdd = append(toAdd, op)
				ret = append(ret, toAdd)
			}
		}
		operatorCache[n] = ret
		return ret
	}
}

func (part1 *Part1) applyOperators(operators []byte) bool {
	result := part1.operands[0]
	operandSum := result
	for i, op := range operators {
		n := part1.operands[i+1]
		if op == '+' {
			result += n
		} else if op == '*' {
			result *= n
		} else {
			concat := strconv.Itoa(result) + strconv.Itoa(n)
			result, _ = strconv.Atoi(concat)
		}
		operandSum += n
	}
	return part1.result == result
}

func main() {
	lines := strings.Split(input, "\n")

	part1 := make([]Part1, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ":")
		result, _ := strconv.Atoi(parts[0])

		parts = strings.Split(parts[1][1:], " ")
		operands := make([]int, len(parts))
		for i, part := range parts {
			operand, _ := strconv.Atoi(part)
			operands[i] = operand
		}

		part1[i] = Part1{result: result, operands: operands}
	}

	fmt.Println(operators(3, []byte("+*")))

	part1Sum := 0
	for _, part := range part1 {
		allOps := operators(len(part.operands)-1, []byte{'+', '*'})
		for _, ops := range allOps {
			if part.applyOperators(ops) {
				part1Sum += part.result
				break
			}
		}
	}
	fmt.Println("part1", part1Sum)

	// Reset memoized function
	operatorCache = make(map[int][][]byte)
	part2Sum := 0
	for _, part := range part1 {
		allOps := operators(len(part.operands)-1, []byte{'+', '*', '|'})
		for _, ops := range allOps {
			if part.applyOperators(ops) {
				part2Sum += part.result
				break
			}
		}
	}
	fmt.Println("part2", part2Sum)
}
