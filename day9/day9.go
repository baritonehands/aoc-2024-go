package main

import (
	_ "embed"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/baritonehands/aoc-2024-go/utils"
	"math"
	"slices"
)

//go:embed input.txt
var input string //= `2333133121414131402`

type File struct {
	id   int
	size uint8
}

func (f *File) Checksum() int {
	return int(f.size) * f.id
}

func free(size uint8) File {
	return File{id: -1, size: size}
}

func isFree(file File) bool {
	return file.id == -1
}

func main() {
	pairs := utils.Partition([]byte(input), 2, 2)
	part1 := make([]File, 0)
	for id, pairSeq := range it.Enumerate(pairs) {
		pair := slices.Collect(pairSeq)
		size := pair[0] - '0'
		part1 = append(part1, slices.Collect(it.Take(it.Repeat(File{id: id, size: 1}), uint(size)))...)
		if len(pair) > 1 {
			part1 = append(part1, slices.Collect(it.Take(it.Repeat(free(1)), uint(pair[1]-'0')))...)
		}
	}

	start := 0
	end := len(part1) - 1
	for start < end {
		for !isFree(part1[start]) {
			start++
		}

		for isFree(part1[end]) {
			end--
		}

		if start >= end {
			break
		}

		part1[start] = part1[end]
		part1[end] = free(1)
	}

	part1Checksum := 0
	for pos, file := range part1 {
		if isFree(file) {
			break
		}
		part1Checksum += pos * file.id
	}

	fmt.Println("part1", part1Checksum)

	filesystem := make([]File, 0)
	freeBlocks := make(map[uint8][]int)
	maxSize := -1
	part2 := make([]int, 0)
	for id, pairSeq := range it.Enumerate(pairs) {
		pair := slices.Collect(pairSeq)
		size := pair[0] - '0'
		filesystem = append(filesystem, File{id: id, size: size})
		part2 = append(part2, slices.Collect(it.Take(it.Repeat(id), uint(size)))...)
		if len(pair) > 1 {
			idx := len(part2)
			freeSize := pair[1] - '0'
			if int(freeSize) > maxSize {
				maxSize = int(freeSize)
			}
			if freeSize > 0 {
				filesystem = append(filesystem, free(freeSize))
				part2 = append(part2, slices.Collect(it.Take(it.Repeat(-1), uint(freeSize)))...)

				if freeBlocks[freeSize] == nil {
					freeBlocks[freeSize] = []int{idx}
				} else {
					freeBlocks[freeSize] = append(freeBlocks[freeSize], idx)
				}
			}
		}
	}
	//fmt.Println("filesystem", filesystem, freeBlocks)

	fsEnd := len(filesystem) - 1
	part2End := len(part2) - int(filesystem[fsEnd].size)
	for fsEnd >= 0 {
		for isFree(filesystem[fsEnd]) {
			fsEnd--
			part2End -= int(filesystem[fsEnd].size)
		}

		if fsEnd < 0 {
			break
		}

		endFile := filesystem[fsEnd]
		leftmostFreeIndexOfSize := math.MaxInt
		var freeSize uint8
		for size := uint8(maxSize); size >= endFile.size; size-- {
			if v, present := freeBlocks[size]; present && len(v) > 0 {
				if v[0] < leftmostFreeIndexOfSize && v[0] < part2End {
					leftmostFreeIndexOfSize = v[0]
					freeSize = size
				}
			}
		}

		if leftmostFreeIndexOfSize != math.MaxInt {
			slices.Replace(part2, leftmostFreeIndexOfSize, leftmostFreeIndexOfSize+int(endFile.size), slices.Repeat([]int{endFile.id}, int(endFile.size))...)
			slices.Replace(part2, part2End, part2End+int(endFile.size), slices.Repeat([]int{-1}, int(endFile.size))...)
			if freeSize > endFile.size {
				rem := freeSize - endFile.size
				freeBlocks[freeSize] = freeBlocks[freeSize][1:]
				freeBlocks[rem] = append(freeBlocks[rem], leftmostFreeIndexOfSize+int(endFile.size))
				slices.Sort(freeBlocks[rem])
			} else {
				freeBlocks[freeSize] = freeBlocks[freeSize][1:]
			}
		}

		fsEnd--
		if fsEnd >= 0 {
			part2End -= int(filesystem[fsEnd].size)
		}
	}
	fmt.Println(part2)

	part2Checksum := 0
	for pos, id := range part2 {
		if id != -1 {
			part2Checksum += pos * id
		}
	}
	fmt.Println("part2", part2Checksum)
}
