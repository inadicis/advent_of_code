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
	rest := amountCycles
	for i := 0; i < amountCycles; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		hash := encode(rows)
		if j, ok := visited[hash]; !foundCycle && ok {
			cycleLength := i - j
			rest = (amountCycles - i) % cycleLength
			foundCycle = true
			// fmt.Printf(" %05d: Found cycle of length %d starting at i: %d. %d steps left\n", i, cycleLength, j, rest)
			// TODO
		}
		if rest <= 0 {
			break
		}
		if !foundCycle {
			visited[encode(rows)] = i
		}
		spinRocks(rows)
		rest--

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

func decode(s string) [][]rune {
	splitLines := strings.Split(s, "\n")
	a := make([][]rune, len(splitLines))
	for i, l := range splitLines {
		r := make([]rune, len(l))
		for j, char := range l {
			r[j] = char
		}
		a[i] = r
	}
	return a
}

func spinRocks(lines [][]rune) {
	// for _, dir := range []Direction{1, 2, 3, 4} {
	// 	fmt.Println(dir)
	// }
	// NORTH
	tiltNorthSouth(lines, false)
	// WEST
	tiltWestEast(lines, false)
	// SOUTH

	tiltNorthSouth(lines, true)
	// EAST
	tiltWestEast(lines, true)
}

func tiltNorthSouth(lines [][]rune, south bool) {

	for col := 0; col < len(lines[0]); col++ {
		initialRow := 0
		lastRow := len(lines) - 1
		step := 1
		if south {
			initialRow, lastRow = lastRow, initialRow
			step = -1
		}
		potentialRow := initialRow
		for row := initialRow; (lastRow-row)*step >= 0; row += step {
			char := lines[row][col]
			if char == '#' {
				potentialRow = row + step
			} else if char == 'O' {
				lines[row][col] = '.'
				lines[potentialRow][col] = 'O'
				potentialRow += step
			} else if char != '.' { // do nothing on points
				panic("unexpected rune")
			}
		}
	}
}
func tiltWestEast(lines [][]rune, east bool) {

	for row := 0; row < len(lines[0]); row++ {
		initialCol := 0
		lastCol := len(lines[0]) - 1
		step := 1
		if east {
			initialCol, lastCol = lastCol, initialCol
			step = -1
		}
		potentialCol := initialCol
		for col := initialCol; (lastCol-col)*step >= 0; col += step {
			char := lines[row][col]
			if char == '#' {
				potentialCol = col + step
			} else if char == 'O' {
				lines[row][col] = '.'
				lines[row][potentialCol] = 'O'
				potentialCol += step
			} else if char != '.' { // do nothing on points
				panic("unexpected rune")
			}
		}
	}
}

func pprint(lines [][]rune) {
	for _, row := range lines {
		for _, char := range row {
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
}

// func tiltRocks(lines [][]rune, upDown bool, reverse bool) {

// 	firstIter := lines
// 	secondIter := lines[0]
// 	if upDown {
// 		firstIter, secondIter := secondIter, firstIter
// 	}

// 	for i := 0; i < len(firstIter); i++ {
// 		potentialPosition := 0
// 		for j := 0; j < len(secondIter); j++ {
// 			if lines[row][col] == '#' {
// 				potentialPosition = row - 1
// 			} else if lines[row][col] == 'O' {
// 				lines[potentialPosition][col] = 'O'
// 				potentialPosition--
// 			} else if lines[row][col] != '.' { // do nothing on points
// 				panic("unexpected rune")
// 			}
// 		}
// 	}

// }

// type Direction int

// const (
// 	NORTH Direction = iota + 1
// 	WEST
// 	SOUTH
// 	EAST
// )

// func tiltRocks(lines []string, direction Direction) {
// 	switch {
// 	case direction == NORTH:

// 	}
// }
