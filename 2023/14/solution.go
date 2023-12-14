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

func same(a [][]rune, b [][]rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, row := range a {
		if len(row) != len(b[i]) {
			return false
		}
		for j, c := range row {
			if c != b[i][j] {
				return false
			}
		}
	}
	return true
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
	// NORTH WEST SOUTH EAST
	for col := 0; col < len(lines[0]); col++ {
		potentialRow := 0
		for row := potentialRow; row < len(lines); row++ {
			char := lines[row][col]
			if char == '#' {
				potentialRow = row + 1
			} else if char == 'O' {
				lines[row][col] = '.'
				lines[potentialRow][col] = 'O'
				potentialRow++
			} else if char != '.' { // do nothing on points
				panic("unexpected rune")
			}
		}
	}
	// fmt.Printf("Rows after tilting North: \n")
	// pprint(lines)
	for row := 0; row < len(lines); row++ {
		// fmt.Printf(" row %02d ... \n", row)
		potentialCol := 0
		for col := potentialCol; col < len(lines[0]); col++ {

			// fmt.Printf("  col %02d ... potential position %d \n", col, potentialPosition)
			char := lines[row][col]
			// fmt.Printf("  char: %q\n", char)
			if char == '#' {
				potentialCol = col + 1
			} else if char == 'O' {
				lines[row][col] = '.'
				lines[row][potentialCol] = 'O'
				potentialCol++
			} else if char != '.' { // do nothing on points
				panic("unexpected rune")
			}
		}
	}
	// fmt.Printf("Rows after tilting west: \n")
	// pprint(lines)
	for col := 0; col < len(lines[0]); col++ {
		potentialRow := len(lines) - 1
		for row := potentialRow; row >= 0; row-- {
			char := lines[row][col]
			if char == '#' {
				potentialRow = row - 1
			} else if char == 'O' {
				lines[row][col] = '.'
				lines[potentialRow][col] = 'O'
				potentialRow--
			} else if char != '.' { // do nothing on points
				panic("unexpected rune")
			}
		}
	}
	// fmt.Printf("Rows after tilting South: \n")
	// pprint(lines)
	for row := 0; row < len(lines); row++ {
		potentialCol := len(lines[0]) - 1
		for col := potentialCol; col >= 0; col-- {
			char := lines[row][col]
			if char == '#' {
				potentialCol = col - 1
			} else if char == 'O' {
				lines[row][col] = '.'
				lines[row][potentialCol] = 'O'
				potentialCol--
			} else if char != '.' { // do nothing on points
				panic("unexpected rune")
			}
		}
	}
	// fmt.Printf("Rows after tilting East: \n")
	// pprint(lines)
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
