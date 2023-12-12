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
		return 0
	}
	if len(amounts) == 0 {
		panic("that should not happen")
	}
	if len(amounts) == 1 {
		// todo check for # . -> actually could stop the recursion at len(amounts) == 0 -> put this in the loop
		if strings.Contains(mask, "#") {

		} else {

			fmt.Printf(pre+" Only 1 amount to place anymore, return %d - %d + 1\n", len(mask), amounts[0])
			return max(0, len(mask)-amounts[0]+1)
		}
	}
	fmt.Printf(pre+"mask %q, amounts: %v", mask, amounts)
	for startIndex := range mask {
		fmt.Printf(pre+" i %02d\n", startIndex)
		endIndex := startIndex + amounts[0]
		if len(mask) < endIndex+2 { //there must be at least one remaining for the next amount
			fmt.Printf(pre + " len(mask) too tiny, break \n")
			break
		}
		workingSprings := mask[:startIndex]

		damagedSprings := mask[startIndex:endIndex]
		workingSprings += string(mask[endIndex])
		remaining := mask[endIndex+1:]
		fmt.Printf(pre+"actual: %q\n", mask)
		w := "................................."
		b := "#################################"
		fmt.Printf(pre+"trying: %q\n", w[:startIndex]+b[startIndex:endIndex]+w[:1]+remaining)
		// fmt.Printf(pre+" working. %q, damaged# %q, remaining %q", workingSprings, damagedSprings, remaining)
		if strings.Contains(workingSprings, "#") {
			fmt.Printf(pre+" working. %q contains '#' -> skipping i %d\n", workingSprings, startIndex)
			continue

		}
		if strings.Contains(damagedSprings, ".") {
			fmt.Printf(pre+" damaged# %q contains '.' -> skipping i %d\n", damagedSprings, startIndex)
			continue
		}
		sum += countPossibilities(remaining, amounts[1:], minNeeded-amounts[0]-1, deepness+1)
		fmt.Printf(pre+"sum increased: %d\n", sum)
	}
	fmt.Printf(pre+"%s %v ==> %d\n", mask, amounts, sum)
	return sum
}
