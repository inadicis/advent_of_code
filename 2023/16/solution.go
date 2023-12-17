package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strings"
)

// func pprint(lines [][]rune) {
// 	for _, row := range lines {
// 		for _, char := range row {
// 			fmt.Printf("%c", char)
// 		}
// 		fmt.Println()
// 	}
// }

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	maze := cleanupData(bytes)
	// pprint(maze)
	result := part2(maze)
	fmt.Printf("\nFinal result: %d", result)

}

// returns a 2D rune matrix (assumes unix-like newlines \n)
func cleanupData(text []byte) [][]rune {
	lines := strings.Split(strings.TrimSpace(string(text)), "\n") // unix files
	rows := make([][]rune, len(lines))
	for i, line := range lines {
		row := make([]rune, len(line))
		for j, char := range line {
			row[j] = char
		}
		rows[i] = row
	}
	return rows
}

// Directions up, down, left or right, defined as two booleans
type ManhattanDirection struct {
	isVertical bool
	sign       bool
	// The "sign" of direction is as if looking at the indexes in the array
	// going right (or down) is positive (true), going left (or up) is going back, is negative (false)
}

var (
	down  = ManhattanDirection{true, true}
	up    = ManhattanDirection{true, false}
	right = ManhattanDirection{false, true}
	left  = ManhattanDirection{false, false}
)

func (d ManhattanDirection) toVector() (v Vector) {
	n := 1
	if !d.sign {
		n = -1
	}
	if d.isVertical {
		v.row = n
	} else {
		v.col = n
	}
	return v
}

type Vector struct {
	row int // (x) 0, -1 or 1 for directions
	col int // (-y)  0, -1 or 1 for directions
}

type Mirror struct {
	isAxis    bool // is horizontal (x-axis) or vertical (y-axis)
	keepsSign bool // for diagonals, if it inverses the direction (positive / negative orientation)
	// keepsSign is "isVertical" for 1-axis mirrors (arbitrarily chosen)
}

func (m Mirror) isVertical() bool {
	// meant for 1-Axis mirrors
	return m.keepsSign
}

var mirrors map[rune]Mirror = map[rune]Mirror{
	'-':  {isAxis: true, keepsSign: false},
	'|':  {isAxis: true, keepsSign: true},
	'/':  {isAxis: false, keepsSign: false},
	'\\': {isAxis: false, keepsSign: true},
}

// Given a light direction and a colliding form, gives next direction(s)
// of the light bean after reflecting.
func reflections(d ManhattanDirection, char rune) []ManhattanDirection {
	outcome := make([]ManhattanDirection, 0, 2)

	mirror, found := mirrors[char]
	if !found { // it's a point, dont't change direction
		return append(outcome, d)
	}

	if !mirror.isAxis { // it's either `/` or `\` (depending on .keepsSign property)
		return append(outcome, ManhattanDirection{!d.isVertical, d.sign == mirror.keepsSign})
		// equivalent to following, which might be easier to understand ?
		// if !mirror.keepsSign {
		// 	outcome = append(outcome, ManhattanDirection{!d.isVertical, !d.sign})
		// } else {
		// 	outcome = append(outcome, ManhattanDirection{!d.isVertical, d.sign})
		// }
	}
	// only `|` or `-` left
	if mirror.isVertical() == d.isVertical {
		// same effect as point, don't change direction
		return append(outcome, d)
	}

	// is a splitter, add two directions
	outcome = append(outcome, ManhattanDirection{!d.isVertical, true})
	return append(outcome, ManhattanDirection{!d.isVertical, false})
}

// Combination of position and direction
type State struct {
	position Vector
	// row       int
	// col       int
	direction ManhattanDirection
}

func getAmountActivated(mirrors [][]rune, startState State) (total int) {
	visited := make(map[State]bool)
	activated := make([][]bool, len(mirrors))
	for i, row := range mirrors {
		activated[i] = make([]bool, len(row))
	}
	queue := list.New()
	queue.PushBack(startState)
	// BFS as it's queue based and not stack based. Both would work.
	for e := queue.Front(); e != nil; e = e.Next() {
		state := e.Value.(State)
		if visited[state] {
			// cycle detected, continue working on the queue without adding new states
			continue
		}
		currentRow := state.position.row
		currentCol := state.position.col
		activated[currentRow][currentCol] = true
		visited[state] = true
		char := mirrors[currentRow][currentCol]
		nextDirections := reflections(state.direction, char)
		for _, d := range nextDirections {
			directionVector := d.toVector()
			row := currentRow + directionVector.row
			col := currentCol + directionVector.col
			if row < 0 || row >= len(mirrors) || col < 0 || col >= len(mirrors[0]) {
				// index out of bounds, light beam going outside the playfield is ignored
				continue
			}
			queue.PushBack(State{position: Vector{row: row, col: col}, direction: d})
		}
	}
	for _, row := range activated {
		for _, active := range row {
			if active {
				total++
			}
		}
	}
	return total
}
func part1(maze [][]rune) (total int) {
	return getAmountActivated(maze, State{Vector{0, 0}, right})
}

func part2(maze [][]rune) (maxTotal int) {
	// how to cache shared ways? 
	for row := 0; row < len(maze); row++ {
		leftToRightTotal := getAmountActivated(maze, State{Vector{row, 0}, right})
		rightToLeftTotal := getAmountActivated(maze, State{Vector{row, len(maze[row]) - 1}, left})
		maxTotal = max(maxTotal, leftToRightTotal, rightToLeftTotal)
	}

	for col := 0; col < len(maze[0]); col++ {
		upToDownTotal := getAmountActivated(maze, State{Vector{0, col}, down})
		downToUpTotal := getAmountActivated(maze, State{Vector{len(maze) - 1, col}, up})
		maxTotal = max(maxTotal, upToDownTotal, downToUpTotal)
	}
	return maxTotal
}
