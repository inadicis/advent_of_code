package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	masks, amounts, err := ExtractData("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(masks)
	// fmt.Println(amounts)
	result := ProcessData(masks, amounts)
	fmt.Printf("\nFinal result: %d", result)
}

func ExtractData(filename string) ([]string, [][]int, error) {
	block, err := os.ReadFile(filename)
	if err != nil {
		return []string{}, [][]int{}, err
	}
	lines := strings.Split(strings.TrimSpace(string(block)), "\r\n")
	rowsMasks := make([]string, len(lines))
	rowsAmounts := make([][]int, len(lines))
	for i, line := range lines {
		maskLine, amountsLine, found := strings.Cut(line, " ")
		if !found {
			return rowsMasks, rowsAmounts, fmt.Errorf("No space found at %d", i)
		}
		amountSplit := strings.Split(amountsLine, ",")
		amounts := make([]int, len(amountSplit))
		for j, amount := range amountSplit {
			a, err := strconv.Atoi(amount)
			if err != nil {
				return rowsMasks, rowsAmounts, err
			}
			amounts[j] = a
		}

		rowsMasks[i] = maskLine
		rowsAmounts[i] = amounts
	}
	return rowsMasks, rowsAmounts, nil
}

func ProcessData(masks []string, amounts [][]int) (total int) {
	for i, mask := range masks {
		minNeeded := -1
		for _, amount := range amounts[i] {
			minNeeded += amount + 1 // at least one more separating from the next one.
		}
		total += countPossibilities(mask, amounts[i], minNeeded)
	}
	return total
}

func countPossibilities(mask string, amounts []int, minNeeded int) (sum int) {
	// ideas
	// - build every possibility based only on len(mask) and on amounts,
	//   then filter out those not-compatible with mask, then count
	// - recursive greadily counting (consuming mask and amounts from left to right)
	// fmt.Printf("%s ||| %v\n", mask, amounts)
	fmt.Printf("countPossibilities(%s, %v, %d)... min>currentLen: %v\n", mask, amounts, minNeeded, minNeeded > len(mask))
	if minNeeded > len(mask) {
		fmt.Printf(" mask is too short -> return 0\n")
		return 0
	}
	if len(amounts) == 0 {
		panic("that should not happen")
	}
	if len(amounts) == 1 {
		fmt.Printf(" Only 1 amount to place anymore, return %d - %d + 1\n", len(mask), amounts[0])
		return max(0, len(mask)-amounts[0]+1)
	}
	for startIndex := range mask {
		fmt.Printf(" %02d: mask %30s, amount: %d\n", startIndex, mask, amounts[0])
		endIndex := startIndex + amounts[0]
		if len(mask) < endIndex+2 { //there must be at least one remaining for the next amount
			return 0
		}
		skipedPart := mask[:startIndex]

		filledPart := mask[startIndex:endIndex]
		skipedPart += string(mask[endIndex])
		restPart := mask[endIndex+1:]
		for _, char := range skipedPart {
			if char == '#' {
				continue
			}
		}
		for _, char := range filledPart {
			if char == '.' {
				continue
			}
		}
		sum += countPossibilities(restPart, amounts[1:], minNeeded-amounts[0]-1)
	}
	fmt.Printf("%s %v ==> %d\n", mask, amounts, sum)
	return sum
}
