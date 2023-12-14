package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nFinal result: %d", part1(string(b)))
}

func part1(input string) int {
	lines := cleanupInput(input)
	currentLoad := 0
	for col := 0; col < len(lines[0]); col++ {
		potentialLoad := len(lines)
		for row := 0; row < len(lines); row++ {
			if lines[row][col] == '#' {
				potentialLoad = len(lines) - row - 1
			} else if lines[row][col] == 'O' {
				currentLoad += potentialLoad
				potentialLoad--
			} else if lines[row][col] != '.' { // do nothing on points
				panic("unexpected rune")
			}
		}
		// fmt.Printf("after col %d: current load: %d\n", col, currentLoad)
		// weight += amountRocks * (amountRocks + 1) / 2
	}
	return currentLoad
}

func part2(input string) int {
	data := cleanupInput(input)
	for _, d := range data {
		fmt.Println(d)
	}
	return -1
}

func cleanupInput(input string) []string {
	return strings.Split(input, "\r\n")
}
