package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

func diag(ws [][]byte, n, ri, ci, rDir, cDir int) string {
	ret := make([]byte, n)
	for di := 0; di < n; di++ {
		r := ri + (di * rDir)
		c := ci + (di * cDir)
		if r >= 0 && c >= 0 && r < len(ws) && c < len(ws[0]) {
			ret[di] = ws[r][c]
		}
	}
	if string(ret) == "XMAS" {
		fmt.Printf("diag (%d,%d) -> %d, %d\n", ri, ci, rDir, cDir)
	}
	return string(ret)
}

func vert(ws [][]byte, n, ri, ci, rDir int) string {
	ret := make([]byte, n)
	for di := 0; di < n; di++ {
		r := ri + (di * rDir)
		if r >= 0 && r < len(ws) {
			ret[di] = ws[r][ci]
		}
	}
	return string(ret)
}

func findWord(ws [][]byte, word string) int {
	cnt := 0
	n := len(word)
	for ri, row := range ws {
		for ci, _ := range row {
			// Right
			if ci <= len(row)-n {
				if string(row[ci:ci+n]) == word {
					fmt.Printf("(%d,%d) -> R\n", ri, ci)
					cnt++
				}
				if diag(ws, n, ri, ci, 1, 1) == word {
					fmt.Printf("(%d,%d) -> DR\n", ri, ci)
					cnt++
				}
				if diag(ws, n, ri, ci, -1, 1) == word {
					fmt.Printf("(%d,%d) -> UR\n", ri, ci)
					cnt++
				}
			}
			if ri <= len(ws)-n && vert(ws, n, ri, ci, 1) == word {
				fmt.Printf("(%d,%d) -> D\n", ri, ci)
				cnt++
			}
			if ci >= n-1 {
				cp := slices.Clone(row[ci-n+1 : ci+1 : ci+1])
				fmt.Println(string(cp))
				slices.Reverse(cp)
				if string(cp) == word {
					fmt.Printf("(%d,%d) -> L\n", ri, ci)
					cnt++
				}
				if diag(ws, n, ri, ci, 1, -1) == word {
					fmt.Printf("(%d,%d) -> DL\n", ri, ci)
					cnt++
				}
				if diag(ws, n, ri, ci, -1, -1) == word {
					fmt.Printf("(%d,%d) -> UL\n", ri, ci)
					cnt++
				}
			}
			if ri >= n-1 && vert(ws, n, ri, ci, -1) == word {
				fmt.Printf("(%d,%d) -> U\n", ri, ci)
				cnt++
			}
		}
	}
	return cnt
}

func main() {
	lines := strings.Split(input, "\n")
	ws := make([][]byte, len(lines))
	for ri, line := range lines {
		ws[ri] = []byte(line)
	}
	println("part1", findWord(ws, "XMAS"))

	part2 := 0
	for ri, row := range ws {
		for ci, _ := range row {
			dr := diag(ws, 3, ri, ci, 1, 1)
			if dr == "MAS" || dr == "SAM" {
				ur := diag(ws, 3, ri+2, ci, -1, 1)
				if ur == "MAS" || ur == "SAM" {
					part2++
				}
			}
		}
	}
	println("part2", part2)
}
