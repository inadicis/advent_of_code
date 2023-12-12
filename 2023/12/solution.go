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

		total += countPossibilities(mask, amounts[i], 0)
	}
	return total
}

func minNeeded(damagedGroupsSizes []int) int {
	m := -1
	for _, amount := range damagedGroupsSizes {
		m += amount + 1 // at least one more separating from the next one.
	}
	return m
}

func countPossibilities(mask string, amounts []int, deepness int) (sum int) {
	// ideas
	// - build every possibility based only on len(mask) and on amounts,
	//   then filter out those not-compatible with mask, then count
	// - recursive greadily counting (consuming mask and amounts from left to right)
	// fmt.Printf("%s ||| %v\n", mask, amounts)
	pre := "|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||"[:deepness]
	fmt.Printf(pre+"countPossibilities(%q, %v, %d)... \n", mask, amounts, minNeeded)
	if len(amounts) == 0 {
		fmt.Printf(pre+" no damagedGroup anymore -> check if # in mask %q to check if valid \n", mask)
		if strings.Contains(mask, "#") {
			fmt.Print(pre + " ==> 0 \n")

			return 0
		}
		fmt.Print(pre + " ==> +1 \n")
		return 1
	}
	m := minNeeded(amounts)
	if m > len(mask) {
		return
	}
	for i := 0; m+i <= len(mask); i++ {
		endIndex := i + amounts[0]
		if endIndex > len(mask) {
			continue
		}
		workingSprings := mask[:i]

		damagedSprings := mask[i:endIndex]
		if len(mask) > endIndex {
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
			// fmt.Printf(pre+" working. %q contains '#' -> skipping i %d\n", workingSprings, i)
			continue

		}
		if strings.Contains(damagedSprings, ".") {
			// fmt.Printf(pre+" damaged# %q contains '.' -> skipping i %d\n", damagedSprings, i)
			continue
		}
		sum += countPossibilities(remaining, amounts[1:], deepness+1)
		// fmt.Printf(pre+"sum increased: %d\n", sum)
	}
	// fmt.Printf(pre+"%s %v ==> %d\n", mask, amounts, sum)
	return sum
}
