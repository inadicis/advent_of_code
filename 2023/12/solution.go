package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	masks, amounts, err := ExtractData("test_data.txt")
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
		total += countPossibilities(mask, amounts[i], minNeeded, 0)
	}
	return total
}

func countPossibilities(mask string, amounts []int, minNeeded int, deepness int) (sum int) {
	// ideas
	// - build every possibility based only on len(mask) and on amounts,
	//   then filter out those not-compatible with mask, then count
	// - recursive greadily counting (consuming mask and amounts from left to right)
	// fmt.Printf("%s ||| %v\n", mask, amounts)
	pre := "|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||"[:deepness]
	fmt.Printf(pre+"countPossibilities(%q, %v, %d)... \n", mask, amounts, minNeeded)
	if minNeeded > len(mask) {
		fmt.Printf(pre + " mask is too short -> return 0\n")
		return
	}
	if len(amounts) == 0 {
		fmt.Print(pre + " COMPLETE: +1\n")
		return 1
	}
	fmt.Printf(pre+"mask %q, amounts: %v, minNeeded: %d\n", mask, amounts, minNeeded)
	for i := 0; minNeeded+i <= len(mask); i++ {
		fmt.Printf(pre+" i %02d\n", i)
		endIndex := i + amounts[0]
		// newMinNeeded := minNeeded - amounts[0] - 1 // - 1 if we need a separator -> not the last one
		if endIndex > len(mask) {
			fmt.Printf(pre+" endIndex %d > len(mask) (%d), break \n", endIndex, len(mask))
			continue
		}
		isLast := len(amounts) == 1
		workingSprings := mask[:i]

		damagedSprings := mask[i:endIndex]
		if !isLast {
			workingSprings += string(mask[endIndex])
		}
		var remaining string
		if len(mask) > endIndex+1 {
			remaining = mask[endIndex+1:]
		}
		// fmt.Printf(pre+"actual: %q\n", mask)
		// w := "................................."
		// b := "#################################"
		// fmt.Printf(pre+"trying: %q\n", w[:startIndex]+b[startIndex:endIndex]+w[:1]+"|"+remaining)
		// fmt.Printf(pre+" working. %q, damaged# %q, remaining %q", workingSprings, damagedSprings, remaining)
		if strings.Contains(workingSprings, "#") {
			fmt.Printf(pre+" working. %q contains '#' -> skipping i %d\n", workingSprings, i)
			continue

		}
		if strings.Contains(damagedSprings, ".") {
			fmt.Printf(pre+" damaged# %q contains '.' -> skipping i %d\n", damagedSprings, i)
			continue
		}
		sum += countPossibilities(remaining, amounts[1:], minNeeded-amounts[0]-1, deepness+1)
		fmt.Printf(pre+"sum increased: %d\n", sum)
	}
	fmt.Printf(pre+"%s %v ==> %d\n", mask, amounts, sum)
	return sum
}
