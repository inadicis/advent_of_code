package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// filename := "test_data.txt"
	filename := "input.txt"
	data, err := ExtractData(filename)

	if err != nil {
		log.Fatal(err)
	}
	result, err := ProcessData(data)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nFinal result: %d", result)
}

func ExtractData(filename string) (patterns [][]string, err error) {
	input, err := os.ReadFile(filename)
	if err != nil {
		return patterns, err
	}
	blocks := strings.Split(strings.TrimSpace(string(input)), "\r\n\r\n")
	patterns = make([][]string, len(blocks))
	for i, block := range blocks {
		patterns[i] = strings.Split(block, "\r\n")
	}

	return patterns, nil
}

func ProcessData(blocks [][]string) (total int, err error) {
	for _, block := range blocks {
		row, col := findSymmetry(block)
		total += col + 100*row
	}
	return total, nil
}

func findSymmetry(block []string) (row int, col int) {
COLAXIS:
	for colAxis := 1; colAxis < len(block[0]); colAxis++ {
		smudgeFound := false
		for _, row := range block {
			distanceToEdge := min(colAxis, len(row)-colAxis)
			for j := 0; j < distanceToEdge; j++ {
				if row[colAxis-j-1] != row[colAxis+j] {
					if smudgeFound {
						continue COLAXIS
					} else {
						smudgeFound = true
					}
					// fmt.Printf("   FALSE -> skip axis %d\n", colAxis)
				}
			}
		}
		if smudgeFound {
			col = colAxis
		}
		// break
	}

ROWAXIS:
	for rowAxis := 1; rowAxis < len(block); rowAxis++ {
		smudgeFound := false
		for i := 0; i < len(block[0]); i++ {
			distanceToEdge := min(rowAxis, len(block)-rowAxis)
			for j := 0; j < distanceToEdge; j++ {
				if block[rowAxis-j-1][i] != block[rowAxis+j][i] {
					if smudgeFound {
						continue ROWAXIS
					} else {
						smudgeFound = true
					}
				}
			}
		}
		if smudgeFound {
			row = rowAxis
		}
		// break
	}
	return row, col
}
