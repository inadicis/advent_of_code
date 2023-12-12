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
	result := ProcessData(masks, amounts, true)
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

func ProcessData(masks []string, amounts [][]int, unfold bool) (total int) {
	if unfold {
		masks, amounts = UnfoldData(masks, amounts)
		fmt.Println("After unfold")
		for i, mask := range masks {
			fmt.Printf("%03d: %q ||| %v\n", i, mask, amounts[i])
		}
	}
	for i, mask := range masks {
		total += countPossibilities(mask, amounts[i])
		fmt.Printf("%03d: currentSum: %d", i, total)
	}
	return total
}

func UnfoldData(masks []string, amounts [][]int) ([]string, [][]int) {
	for i, mask := range masks {
		masks[i] = strings.Join([]string{mask, mask, mask, mask, mask}, "?")
		newAmounts := make([]int, len(amounts[i])*5)
		for j := 0; j < 5; j++ {
			for k, a := range amounts[i] {
				newAmounts[j*len(amounts[i])+k] = a
			}
		}
		amounts[i] = newAmounts
	}
	return masks, amounts
}

func minNeeded(damagedGroupsSizes []int) int {
	m := -1
	for _, amount := range damagedGroupsSizes {
		m += amount + 1 // at least one more separating from the next one.
	}
	return m
}

func countPossibilities(mask string, amounts []int) (sum int) {
	// recursive greadily counting (consuming mask and amounts from left to right)
	if len(amounts) == 0 {
		if strings.Contains(mask, "#") {
			return 0
		}
		return 1
	}
	min := minNeeded(amounts)
	if min > len(mask) {
		return
	}
	for offset := 0; min+offset <= len(mask); offset++ {
		endIndex := offset + amounts[0]
		workingSprings := mask[:offset]
		damagedSprings := mask[offset:endIndex]
		var remainingSprings string

		if len(mask) > endIndex {
			workingSprings += string(mask[endIndex])
			if len(mask) > endIndex+1 {
				remainingSprings = mask[endIndex+1:]
			}

		}

		if strings.Contains(workingSprings, "#") || strings.Contains(damagedSprings, ".") {
			continue

		}
		sum += countPossibilities(remainingSprings, amounts[1:])
	}
	return sum
}
