package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
)

func main() {
	rows, err := ExtractData("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(rows)
	sumForward := 0
	sumBackward := 0

	for _, row := range rows {
		nextForward, err := ExtrapolateForward(row)
		if err != nil {
			log.Fatal(err)
		}
		sumForward += nextForward

		nextBackward, err := ExtrapolateBackward(row)
		if err != nil {
			log.Fatal(err)
		}
		sumBackward += nextBackward

	}
	fmt.Println(sumForward, sumBackward) // 952946799

	// fmt.Println(quote.Glass())
	// fmt.Println(combin.Binomial(4, 4))
}

func ExtrapolateForward(row []int) (int, error) {
	deepness, found := FindDeepnessOf0Row(row)
	if !found {
		return 0, fmt.Errorf("No 0-filled row for start row %v", row)
	}
	deeperRowLastValue := 0
	for d := deepness - 1; d >= 0; d-- {
		currentRowLastValue := calcValue(row, d, len(row)-d-1)
		deeperRowLastValue = deeperRowLastValue + currentRowLastValue
	}
	return deeperRowLastValue, nil

}

func ExtrapolateBackward(row []int) (int, error) {
	deepness, found := FindDeepnessOf0Row(row)
	if !found {
		return 0, fmt.Errorf("No 0-filled row for start row %v", row)
	}
	deeperRowFirstValue := 0
	for d := deepness - 1; d >= 0; d-- {
		currentRowFirstValue := calcValue(row, d, 0)
		// we already calculated this one, could cache it

		deeperRowFirstValue = currentRowFirstValue - deeperRowFirstValue
	}
	return deeperRowFirstValue, nil

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
		value := combin.Binomial(deepness, k) * firstRow[index+deepness-k]
		if sign {
			result += value
		} else {
			result -= value
		}
	}
	return result
}
