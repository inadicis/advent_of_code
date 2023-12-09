package main

import (
	"os"
	"strconv"
	"strings"
)

func main() {
	// expectedResult := [][]int{
	// 	{0, 3, 6, 9, 12, 15},
	// 	{1, 3, 6, 10, 15, 21},
	// 	{1, 0, 13, 16, 21, 30, 45},
	// 	{55, -95, -349, -669, -918, -793, 264, 3180},
	// }
}

func ExtractData(filename string) ([][]int, error) {
	block, err := os.ReadFile(filename)
	if err != nil {
		return [][]int{}, err
	}
	lines := strings.Split(strings.TrimSpace(string(block)), "\r\n")
	rows := make([][]int, len(lines))
	for i, line := range lines {
		numbersLine := strings.Split(line, " ")
		newRow := make([]int, len(numbersLine))
		for j, number := range numbersLine {
			n, err := strconv.Atoi(number)
			if err != nil {
				return rows, err
			}
			newRow[j] = n
		}
		rows[i] = newRow
	}
	return rows, nil
}
