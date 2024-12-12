package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/baritonehands/aoc-2024-go/utils"
	"maps"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

type Region struct {
	plant                    byte
	points                   map[utils.Point]bool
	perimeter                int
	top, right, bottom, left []utils.Point
}

func (r Region) String() string {
	return fmt.Sprintf("%v: perim=%d, sides=%v, %v", string(r.plant), r.perimeter, r.Sides(), slices.Collect(maps.Keys(r.points)))
}

func (r *Region) Sides() string {
	return fmt.Sprintf("|t=%v,r=%v,b=%v,l=%v|", r.top, r.right, r.bottom, r.left)
}

func (r *Region) Area() int {
	return len(r.points)
}

func part2Partition(result [][]utils.Point, side utils.Point) [][]utils.Point {
	if len(result) == 0 {
		result = append(result, []utils.Point{side})
	} else {
		lastResult := result[len(result)-1]
		last := lastResult[len(lastResult)-1]

		if (last.X == side.X && side.Y-last.Y == 1) || (last.Y == side.Y && side.X-last.X == 1) {
			result[len(result)-1] = append(lastResult, side)
		} else {
			result = append(result, []utils.Point{side})
		}
	}
	return result
}

func main() {
	lines := strings.Split(input, "\n")
	yMax := len(lines) - 1
	xMax := len(lines[0]) - 1
	seen := map[utils.Point]bool{}
	regions := []Region{}
	for ri, row := range lines {
		for ci, col := range row {
			point := utils.Point{X: ci, Y: ri}

			if _, ok := seen[point]; !ok {
				curRegion := Region{plant: byte(col), points: map[utils.Point]bool{}}

				regionBfs := []utils.Point{point}
				for len(regionBfs) > 0 {
					cur := regionBfs[0]
					regionBfs = regionBfs[1:]
					if seen[cur] {
						continue
					}

					curRegion.points[cur] = true
					seen[cur] = true

					if cur.X == 0 || cur.X == xMax {
						if cur.X == 0 {
							curRegion.left = append(curRegion.left, cur)
						} else if cur.X == xMax {
							curRegion.right = append(curRegion.right, cur)
						}
						curRegion.perimeter++
					}

					if cur.Y == 0 || cur.Y == yMax {
						if cur.Y == 0 {
							curRegion.top = append(curRegion.top, cur)
						} else if cur.Y == yMax {
							curRegion.bottom = append(curRegion.bottom, cur)
						}
						curRegion.perimeter++
					}

					neighbors := cur.OrthogonalNeighbors(xMax, yMax)
					for _, neighbor := range neighbors {
						if lines[neighbor.Y][neighbor.X] != uint8(col) {
							curRegion.perimeter++
							if neighbor.Y > cur.Y {
								curRegion.bottom = append(curRegion.bottom, cur)
							}
							if neighbor.Y < cur.Y {
								curRegion.top = append(curRegion.top, cur)
							}
							if neighbor.X > cur.X {
								curRegion.right = append(curRegion.right, cur)
							}
							if neighbor.X < cur.X {
								curRegion.left = append(curRegion.left, cur)
							}
						}
						if !seen[neighbor] && lines[neighbor.Y][neighbor.X] == uint8(col) {
							regionBfs = append(regionBfs, neighbor)
						}
					}
				}

				regions = append(regions, curRegion)
			}
		}
	}

	part1 := 0
	for _, region := range regions {
		price := region.Area() * region.perimeter
		part1 += price
		//fmt.Printf("%v = %d\n", region, price)
	}
	fmt.Println("part1", part1)

	part2 := 0
	for _, region := range regions {
		slices.SortFunc(region.top, utils.PointCompareYX)
		slices.SortFunc(region.right, utils.PointCompareXY)
		slices.SortFunc(region.bottom, utils.PointCompareYX)
		slices.SortFunc(region.left, utils.PointCompareXY)

		top := it.Fold(slices.Values(region.top), part2Partition, [][]utils.Point{})
		right := it.Fold(slices.Values(region.right), part2Partition, [][]utils.Point{})
		bottom := it.Fold(slices.Values(region.bottom), part2Partition, [][]utils.Point{})
		left := it.Fold(slices.Values(region.left), part2Partition, [][]utils.Point{})
		//fmt.Printf("%v|t=%v,r=%v,b=%v,l=%v|\n", string(region.plant), top, right, bottom, left)

		sides := 0
		sides += len(top)
		sides += len(right)
		sides += len(bottom)
		sides += len(left)

		//fmt.Printf("%v: %d * %d = %d\n", string(region.plant), region.Area(), sides, region.Area()*sides)
		price := region.Area() * sides
		part2 += price
		//fmt.Printf("%v = %d\n", region, price)
	}
	fmt.Println("part2", part2)
}
