package main

import (
	_ "embed"
	"fmt"
	"github.com/baritonehands/aoc-2024-go/utils"
	"maps"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

//= `kh-tc
//qp-kh
//de-cg
//ka-co
//yn-aq
//qp-ub
//cg-tb
//vc-aq
//tb-ka
//wh-tc
//yn-cg
//kh-ub
//ta-co
//de-co
//tc-td
//tb-wq
//wh-td
//ta-ka
//td-qp
//aq-cg
//wq-ub
//ub-vc
//de-ta
//wq-aq
//wq-vc
//wh-yn
//ka-de
//kh-ta
//co-tc
//wh-qp
//tb-vc
//td-yn`

type Graph map[string][]string

type Triple struct {
	v1, v2, v3 string
}

func main() {
	lines := strings.Split(input, "\n")
	fmt.Println(lines)

	graph := Graph{}
	for _, l := range lines {
		parts := strings.Split(l, "-")
		graph[parts[0]] = append(graph[parts[0]], parts[1])
		graph[parts[1]] = append(graph[parts[1]], parts[0])
	}
	fmt.Println(graph)

	triples := map[Triple]bool{}
	for first, firstConns := range graph {
		for _, second := range firstConns {
			if second != first {
				secondConns := graph[second]
				for _, third := range secondConns {
					if third != first && third != second && slices.Contains(firstConns, third) {
						vs := []string{first, second, third}
						slices.Sort(vs)

						triple := Triple{vs[0], vs[1], vs[2]}
						triples[triple] = true
					}
				}
			}
		}
	}
	fmt.Println(len(triples), triples)

	part1 := 0
	for triple, _ := range triples {
		if triple.v1[0] == 't' || triple.v2[0] == 't' || triple.v3[0] == 't' {
			part1++
		}
	}
	fmt.Println("part1", part1)

	part2 := map[string]bool{}
	for first, firstConns := range graph {
		toCheck := utils.SeqSet(slices.Values(firstConns))
		toCheck[first] = true
		//opts := pq.NewQueue(func(set map[string]bool) int {
		//	return -len(set)
		//})

		for _, other := range firstConns {
			if toCheck[other] {
				otherSet := utils.SeqSet(slices.Values(graph[other]))
				otherSet[other] = true
				toCheck = utils.SetIntersection(toCheck, otherSet)
			}
		}
		if len(toCheck) > len(part2) {
			part2 = toCheck
		}
	}
	fmt.Println(part2)

	sorted := slices.Collect(maps.Keys(part2))
	slices.Sort(sorted)
	fmt.Println(strings.Join(sorted, ","))
}
