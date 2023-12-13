package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filename := "test_data.txt"
	// filename := "input.txt"
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

func ExtractData(filename string) (lines []string, err error) {
	block, err := os.ReadFile(filename)
	if err != nil {
		return lines, err
	}
	lines = strings.Split(strings.TrimSpace(string(block)), "\r\n")
	// rowsMasks := make([]string, len(lines))
	// rowsAmounts := make([][]int, len(lines))
	// for i, line := range lines {
	// 	maskLine, amountsLine, found := strings.Cut(line, " ")
	// 	if !found {
	// 		return rowsMasks, rowsAmounts, fmt.Errorf("No space found at %d", i)
	// 	}
	// 	amountSplit := strings.Split(amountsLine, ",")
	// 	amounts := make([]int, len(amountSplit))
	// 	for j, amount := range amountSplit {
	// 		a, err := strconv.Atoi(amount)
	// 		if err != nil {
	// 			return rowsMasks, rowsAmounts, err
	// 		}
	// 		amounts[j] = a
	// 	}

	// 	rowsMasks[i] = maskLine
	// 	rowsAmounts[i] = amounts
	// }
	return lines, nil
}

func ProcessData(lines []string) (total int, err error) {

	return total, nil
}
