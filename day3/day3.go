package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	pattern, _ := regexp.Compile(`mul\((\d+),(\d+)\)|do\(\)|don't\(\)`)
	result := pattern.FindAllStringSubmatch(input, -1)
	part1 := 0
	part2 := 0

	enabled := true
	for _, entry := range result {
		//fmt.Printf("%v", entry)
		if strings.HasPrefix(entry[0], "mul") {
			l, _ := strconv.Atoi(entry[1])
			r, _ := strconv.Atoi(entry[2])
			part1 += l * r

			if enabled {
				part2 += l * r
			}
		} else if strings.HasPrefix(entry[0], "don't") {
			enabled = false
		} else {
			enabled = true
		}

	}
	fmt.Println("part1", part1)
	fmt.Println("part2", part2)

}
