package main

import (
	_ "embed"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//	`Register A: 729
//Register B: 0
//Register C: 0
//
//Program: 0,1,5,4,3,0`

var debug = false

type Computer struct {
	IP      int
	A, B, C int
	output  []string
}

func (state *Computer) HexString() string {
	return fmt.Sprintf("A:%x B:%x C:%x, %v", state.A, state.B, state.C, state.output)
}

func (state *Computer) comboOperand(operand int) int {
	if operand >= 0 && operand <= 3 {
		return operand
	} else if operand == 4 {
		return state.A
	} else if operand == 5 {
		return state.B
	} else if operand == 6 {
		return state.C
	}
	panic("Shouldn't reach here")
}

func processInstruction(state *Computer, program []string) {
	opcode, _ := strconv.Atoi(program[state.IP])
	literal, _ := strconv.Atoi(program[state.IP+1])
	combo := state.comboOperand(literal)

	incIP := true

	if debug {
		fmt.Printf("Op %d [%d,%d] ", opcode, literal, combo)
	}
	switch opcode {
	case 0:
		state.adv(combo)
	case 1:
		state.bxl(literal)
	case 2:
		state.bst(combo)
	case 3:
		incIP = !state.jnz(literal)
	case 4:
		state.bxc()
	case 5:
		state.out(combo)
	case 6:
		state.bdv(combo)
	case 7:
		state.cdv(combo)
	}

	if incIP {
		state.IP += 2
	}
}

func (state *Computer) adv(denominator int) {
	state.A = state.A / int(math.Pow(2.0, float64(denominator)))
}

func (state *Computer) bxl(literal int) {
	state.B = state.B ^ literal
}

func (state *Computer) bst(combo int) {
	state.B = combo % 8
}

func (state *Computer) jnz(literal int) bool {
	if state.A == 0 {
		return false
	}
	state.IP = literal
	return true
}

func (state *Computer) bxc() {
	state.B = state.B ^ state.C
}

func (state *Computer) out(combo int) {
	state.output = append(state.output, fmt.Sprint(combo%8))
}

func (state *Computer) bdv(denominator int) {
	state.B = state.A / int(math.Pow(2.0, float64(denominator)))
}

func (state *Computer) cdv(denominator int) {
	state.C = state.A / int(math.Pow(2.0, float64(denominator)))
}

func main() {
	lines := strings.Split(input, "\n")
	state := Computer{}
	state.A, _ = strconv.Atoi(lines[0][12:])
	state.B, _ = strconv.Atoi(lines[1][12:])
	state.C, _ = strconv.Atoi(lines[2][12:])

	program := strings.Split(lines[4][9:], ",")

	initialState := state

	for state.IP < len(program)-1 {
		processInstruction(&state, program)
		fmt.Println(state.HexString())
	}
	fmt.Println(state, program)
	fmt.Println("part1", strings.Join(state.output, ","))

	i := 1
	for {
		state = initialState
		state.A = i
		if i%100000 == 0 {
			fmt.Printf("%d start\n", i)
		}
		for state.IP < len(program)-1 {
			processInstruction(&state, program)
		}
		if slices.Equal(state.output, program) {
			fmt.Println("part2", i)
			break
		} else if len(program) > len(state.output) {
			i *= 2
		} else if len(program) == len(state.output) {
			for j := len(program) - 1; j >= 0; j-- {
				if program[j] != state.output[j] {
					i += int(math.Pow(8, float64(j)))
					break
				}
			}
		} else {
			i /= 2
		}

		if i%100000 == 0 {
			fmt.Printf("%d: %v\n", i, state)
		}

	}
}
