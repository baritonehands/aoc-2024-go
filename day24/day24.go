package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"maps"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//= `x00: 1
//x01: 0
//x02: 1
//x03: 1
//x04: 0
//y00: 1
//y01: 1
//y02: 1
//y03: 1
//y04: 1`

//go:embed input2.txt
var input2 string

//= `ntg XOR fgs -> mjb
//y02 OR x01 -> tnw
//kwq OR kpj -> z05
//x00 OR x03 -> fst
//tgd XOR rvg -> z01
//vdt OR tnw -> bfw
//bfw AND frj -> z10
//ffh OR nrd -> bqk
//y00 AND y03 -> djm
//y03 OR y00 -> psh
//bqk OR frj -> z08
//tnw OR fst -> frj
//gnj AND tgd -> z11
//bfw XOR mjb -> z00
//x03 OR x00 -> vdt
//gnj AND wpb -> z02
//x04 AND y00 -> kjc
//djm OR pbm -> qhw
//nrd AND vdt -> hwm
//kjc AND fst -> rvg
//y04 OR y02 -> fgs
//y01 AND x02 -> pbm
//ntg OR kjc -> kwq
//psh XOR fgs -> tgd
//qhw XOR tgd -> z09
//pbm OR djm -> kpj
//x03 XOR y03 -> ffh
//x00 XOR y04 -> ntg
//bfw OR bqk -> z06
//nrd XOR fgs -> wpb
//frj XOR qhw -> z04
//bqk OR frj -> z07
//y03 OR x01 -> nrd
//hwm AND bqk -> z03
//tgd XOR rvg -> z12
//tnw OR pbm -> gnj`

type Op struct {
	output             string
	operator           string
	operand1, operand2 string
}

func (op Op) Perform(state map[string]bool) bool {
	var l, r, found bool
	if l, found = state[op.operand1]; !found {
		panic("No left operand")
	}
	if r, found = state[op.operand2]; !found {
		panic("No right operand")
	}

	switch op.operator {
	case "AND":
		return l && r
	case "OR":
		return l || r
	case "XOR":
		if l && !r || !l && r {
			return true
		} else {
			return false
		}
	default:
		panic("Unknown operator")
	}
}

type DirectedGraph[T comparable] struct {
	vertices      map[T][]T
	incomingEdges map[T][]T
}

func NewDirectedGraph[T comparable]() DirectedGraph[T] {
	return DirectedGraph[T]{vertices: make(map[T][]T), incomingEdges: make(map[T][]T)}
}

func (graph *DirectedGraph[T]) AddVertex(label T) {
	if graph.vertices[label] == nil {
		graph.vertices[label] = []T{}
	}
	if graph.incomingEdges[label] == nil {
		graph.incomingEdges[label] = []T{}
	}
}

func (graph *DirectedGraph[T]) AddEdge(label1, label2 T) {
	graph.vertices[label1] = append(graph.vertices[label1], label2)
	graph.incomingEdges[label2] = append(graph.incomingEdges[label2], label1)
}

func (graph *DirectedGraph[T]) RemoveEdge(label1, label2 T) {
	eV1 := graph.vertices[label1]
	eV2 := graph.incomingEdges[label2]
	if eV1 != nil {
		idx := slices.Index(eV1, label2)
		if idx != -1 {
			graph.vertices[label1] = append(eV1[:idx], eV1[idx+1:]...)
		}
	}
	if eV2 != nil {
		idx := slices.Index(eV2, label1)
		if idx != -1 {
			graph.incomingEdges[label2] = append(eV2[:idx], eV2[idx+1:]...)
		}
	}
}

func (graph *DirectedGraph[T]) TopologicalSort() []T {
	ret := []T{}
	s := []T{}

	for k, nodes := range graph.incomingEdges {
		if len(nodes) == 0 {
			s = append(s, k)
		}
	}

	for len(s) > 0 {
		n := s[0]
		s = s[1:]
		ret = append(ret, n)

		edges := slices.Clone(graph.vertices[n])
		for _, m := range edges {
			graph.RemoveEdge(n, m)

			if len(graph.incomingEdges[m]) == 0 {
				s = append(s, m)
			}
		}
	}

	allEmpty := it.All(it.Map(maps.Values(graph.incomingEdges), func(list []T) bool {
		return len(list) == 0
	}))

	if !allEmpty {
		panic("Graph has at least one allEmpty")
	}

	return ret
}

func bitsToInt(result map[string]bool, prefix string) int {
	value := 0
	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("%s%02d", prefix, i)
		if result[k] {
			value |= 1 << i
		}
	}
	return value
}

func main() {
	lines := strings.Split(input, "\n")
	inputs := make(map[string]bool, len(lines))
	for _, l := range lines {
		parts := strings.Split(l, ": ")
		inputs[parts[0]], _ = strconv.ParseBool(parts[1])
	}
	fmt.Println(inputs)

	lines2 := strings.Split(input2, "\n")
	ops := make(map[string]Op, len(lines2))
	opsByOperand := make(map[string][]Op, len(lines2))
	for _, l := range lines2 {
		parts := strings.Split(l, " -> ")
		opParts := strings.Split(parts[0], " ")
		name := parts[1]
		op := Op{operator: opParts[1], operand1: opParts[0], operand2: opParts[2], output: name}
		ops[name] = op
		opsByOperand[op.operand1] = append(opsByOperand[op.operand1], op)
		opsByOperand[op.operand2] = append(opsByOperand[op.operand2], op)
	}
	fmt.Println(ops)
	fmt.Println(opsByOperand)

	graph := NewDirectedGraph[string]()
	for _, op := range ops {
		graph.AddVertex(op.operand1)
		graph.AddVertex(op.operand2)
		graph.AddVertex(op.output)
		graph.AddEdge(op.operand1, op.output)
		graph.AddEdge(op.operand2, op.output)
	}

	processingOrder := graph.TopologicalSort()
	fmt.Println(processingOrder)
	result := maps.Clone(inputs)
	for _, v := range processingOrder {
		if op, found := ops[v]; found {
			result[op.output] = op.Perform(result)
		}
	}
	fmt.Println(result)

	z := bitsToInt(result, "z")
	fmt.Printf("%b: %d\n", z, z)

	x := bitsToInt(result, "x")
	y := bitsToInt(result, "y")

	fmt.Printf("%d + %d = %d vs %d\n", x, y, x+y, z)
	wrong := 0
	for i := 0; i < 100; i++ {
		bit := 1 << i
		if ((x & bit) & (y & bit)) != (z & bit) {
			op1Name := fmt.Sprintf("x%02d", i)
			op2Name := fmt.Sprintf("y%02d", i)
			fmt.Printf("%v <-> %v\n", opsByOperand[op1Name], opsByOperand[op2Name])
			wrong |= bit
		}
	}
	fmt.Printf("%b\n", wrong)
}
