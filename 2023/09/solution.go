package main

import (
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
)

func main() {
	rows := ExtractData("input.txt")

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

func FindDeepnessOf0Row(row []int) (d int, found bool) {
DEEP:
	for d = 0; d < len(row); d++ {
		for i := 0; i < len(row)-d; i++ {
			if calcValue(row, d, i) != 0 {
				continue DEEP
			}
		}
		return d, true
	}
	return
}

func calcValue(firstRow []int, deepness int, index int) int {
	result := 0
	for k := 0; k <= deepness; k++ {
		sign := k%2 == 0
		value := combin.Binomial(k, deepness) * firstRow[index+deepness-k]
		if sign {
			result += value
		} else {
			result -= value
		}
	}
	return result
}
