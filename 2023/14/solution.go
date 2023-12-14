package main

import (
	"bytes"
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
	fmt.Printf("\nFinal result: %d", part2(string(b)))
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

func cleanupInput(input string) []string {
	return strings.Split(input, "\r\n")
}

func part2(input string) int {
	rows := StringsToRunes(cleanupInput(input))
	visited := make(map[string]int)

	amountCycles := 1000000000
	foundCycle := false
	for i := amountCycles; i > 0; i-- {
		hash := encode(rows)
		if !foundCycle {
			if j, ok := visited[hash]; ok {
				cycleLength := i - j
				i = i % cycleLength
				foundCycle = true
			} else {
				visited[hash] = i
			}
		}
		spinRocks(rows)
	}
	return countWeight(rows)
}

func countWeight(rows [][]rune) (weight int) {
	for i, row := range rows {
		for _, c := range row {
			if c == 'O' {
				weight += len(rows) - i
			}
		}
	}
	return weight
}

func StringToRunes(s string) []rune {
	r := make([]rune, len(s))
	for i, c := range s {
		r[i] = c
	}
	return r
}

func StringsToRunes(a []string) [][]rune {
	rows := make([][]rune, len(a))
	for i, s := range a {
		rows[i] = StringToRunes(s)
	}
	return rows
}

func encode(a [][]rune) string {
	var buffer bytes.Buffer
	for _, row := range a {
		for _, char := range row {
			buffer.WriteRune(char)
		}
		buffer.WriteRune('\n')
	}
	return buffer.String()
}

func spinRocks(lines [][]rune) {
	// NORTH WEST SOUTH EAST
	tilt(lines, true, false)
	tilt(lines, false, false)
	tilt(lines, true, true)
	tilt(lines, false, true)
}

func tilt(lines [][]rune, upDown bool, reverse bool) {
	// generic tilt, but harder to understand
	firstIter := len(lines[0])
	secondIter := len(lines)
	if !upDown {
		firstIter, secondIter = secondIter, firstIter
	}
	for i := 0; i < firstIter; i++ {
		initialJ := 0
		lastJ := secondIter - 1
		step := 1
		if reverse {
			initialJ, lastJ = lastJ, initialJ
			step = -1
		}
		potentialJ := initialJ
		for j := initialJ; (lastJ-j)*step >= 0; j += step {
			row, col := j, i
			potentialRow, potentialCol := potentialJ, col
			if !upDown {
				row, col = col, row
				potentialRow, potentialCol = row, potentialJ
			}
			char := lines[row][col]
			switch {
			case char == '#':
				potentialJ = j + step
			case char == 'O':
				lines[row][col] = '.'
				lines[potentialRow][potentialCol] = 'O'
				potentialJ += step
			}
		}
	}
}
