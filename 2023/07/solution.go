package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	r, err := getResult("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	// s := "Hello World"
	// for i, r := range s {
	// 	fmt.Printf("Range %d, type:%T, r= %#v, r=%d \n", i, r, r, r)
	// }
	// for i := 0; i < len(s); i++ {
	// 	fmt.Printf("index accessing hello[%d], type: %T v=%#v, r=%d \n", i, s[i], s[i], s[i])
	// }

	fmt.Println(r)
}

type SortByHand []string

func (a SortByHand) Len() int           { return len(a) }
func (a SortByHand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByHand) Less(i, j int) bool { return a[i] < a[j] }

func compareHands(line1 string, line2 string) int {
	hand1, _, _ := strings.Cut(line1, " ")
	hand2, _, _ := strings.Cut(line2, " ")
	primary1, secondary1 := getHandValue(hand1, 'J')
	primary2, secondary2 := getHandValue(hand2, 'J')
	primaryDiff := primary1 - primary2
	if primaryDiff != 0 {
		return primaryDiff
	}
	secondaryDiff := secondary1 - secondary2
	if secondaryDiff != 0 {
		return secondaryDiff
	}
	cardValues := map[rune]int{
		'A': 14,
		'K': 13,
		'Q': 12,
		'J': 0,
		'T': 10,
	}
	for i, card1 := range hand1 {
		card2 := hand2[i]
		cardDiff := getCardValue(card1, cardValues) - getCardValue(rune(card2), cardValues)
		// GOLANG? why iterate over a string provides runes (int32), but accessing an index provides a byte (uint8)?
		if cardDiff != 0 {
			return cardDiff
		}
	}
	return 0
}

func getCardValue(card rune, cardValues map[rune]int) int {
	if unicode.IsDigit(card) {
		return int(card - '0')
	}
	return cardValues[card]
}

func getHandValue(hand string, joker rune) (int, int) {
	counts := map[rune]int{}
	for _, card := range hand {
		counts[card]++
	}
	primary, secondary := 0, 0
	for card, count := range counts {
		if card != joker {
			if count > primary {
				secondary = primary
				primary = count
			} else if count > secondary {
				secondary = count
			}
		}
	}
	return primary + counts[joker], secondary

}

func getResult(filename string) (int, error) {
	text, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	handsLines := strings.Split(strings.TrimSpace(string(text)), "\r\n")
	slices.SortFunc(handsLines, compareHands)

	total := 0
	for i, line := range handsLines {
		_, bidStr, _ := strings.Cut(line, " ")
		bid, err := strconv.Atoi(bidStr)
		if err != nil {
			return 0, err
		}
		total += bid * (i + 1)
	}

	return total, nil
}
