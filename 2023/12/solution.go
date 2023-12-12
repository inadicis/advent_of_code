package main

import (
	"bytes"
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
	}
	cache := map[Input]int{}
	for i, mask := range masks {
		var additional int
		additional, cache = countPossibilities(Input{mask, encode(amounts[i])}, cache)
		total += additional
		// fmt.Printf("%03d: currentSum: %d, cache size: %d\n", i, total, len(cache))
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

type Input struct {
	mask    string
	amounts string
}

func encode(amounts []int) string {
	var buffer bytes.Buffer
	for i, a := range amounts {
		if i != 0 {
			buffer.WriteRune(',')
		}
		buffer.WriteString(fmt.Sprintf("%d", a))
	}
	s := buffer.String()
	return s
}

func decode(s string) []int {
	split := strings.Split(s, ",")
	numbers := make([]int, len(split))
	for i, n := range split {
		number, err := strconv.Atoi(n)
		if err != nil {
			panic("Number not recognized")
		}
		numbers[i] = number
	}
	return numbers
}

func countPossibilities(in Input, cache map[Input]int) (sum int, updatedCache map[Input]int) {
	// recursive greadily counting (consuming mask and amounts from left to right)
	if a, ok := cache[in]; ok {
		return a, cache
	}
	if len(in.amounts) == 0 {
		if strings.Contains(in.mask, "#") {
			cache[in] = 0
			return 0, cache
		}
		cache[in] = 1
		return 1, cache
	}
	amounts := decode(in.amounts)
	min := minNeeded(amounts)
	if min > len(in.mask) {
		return
	}
	for offset := 0; min+offset <= len(in.mask); offset++ {
		endIndex := offset + amounts[0]
		workingSprings := in.mask[:offset]
		damagedSprings := in.mask[offset:endIndex]
		var remainingSprings string

		if len(in.mask) > endIndex {
			workingSprings += string(in.mask[endIndex])
			if len(in.mask) > endIndex+1 {
				remainingSprings = in.mask[endIndex+1:]
			}

		}

		if strings.Contains(workingSprings, "#") || strings.Contains(damagedSprings, ".") {
			continue

		}
		var additionalPossibilities int
		additionalPossibilities, cache = countPossibilities(Input{remainingSprings, encode(amounts[1:])}, cache)
		sum += additionalPossibilities
	}
	cache[in] = sum
	return sum, cache
}
