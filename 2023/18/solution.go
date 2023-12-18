package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func pprint(lines [][]bool) {
	for _, row := range lines {
		for _, isHole := range row {
			if isHole {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	bytes, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	instructions := cleanupData(bytes)
	result := part2(instructions)
	fmt.Printf("\nFinal result: %d", result)
}

type Instruction struct {
	direction ManhattanDirection
	amount    uint64
	color     string
}

var directionMap map[string]ManhattanDirection = map[string]ManhattanDirection{
	"R": right,
	"L": left,
	"U": up,
	"D": down,
}

// returns a 2D rune matrix (assumes unix-like newlines \n)
func cleanupData(text []byte) []Instruction {
	lines := strings.Split(strings.TrimSpace(string(text)), "\n") // unix files
	instructions := make([]Instruction, len(lines))
	for i, line := range lines {
		elements := strings.Split(line, " ")
		d := directionMap[elements[0]]
		a, err := strconv.Atoi(elements[1])
		if err != nil {
			panic(err)
		}
		instructions[i] = Instruction{
			direction: d,
			amount:    uint64(a),
			color:     elements[2],
		}
	}
	return instructions
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

func (v Vector) add(v2 Vector) Vector {
	return Vector{v.row + v2.row, v.col + v2.col}
}

func (v Vector) isOutOfBounds(maze [][]int) bool {
	return v.row < 0 || v.col < 0 || v.row >= len(maze) || v.col >= len(maze[v.row])
}

func getInitialState(instructions []Instruction) (dimensions Vector, startPosition Vector) {
	currentPos := Vector{}

	var minRow, minCol, maxRow, maxCol int
	for _, instruction := range instructions {
		for i := uint64(0); i < instruction.amount; i++ {
			currentPos = currentPos.add(instruction.direction.toVector())
		}
		maxRow = max(currentPos.row, maxRow)
		maxCol = max(currentPos.col, maxCol)
		minRow = min(currentPos.row, minRow) // offset to startPos
		minCol = min(currentPos.col, minCol)
	}
	startPosition = Vector{-minRow, -minCol}
	dimensions = Vector{maxRow - minRow + 1, maxCol - minCol + 1}
	return dimensions, startPosition
}

func dig(instructions []Instruction) [][]bool {
	dimensions, startPos := getInitialState(instructions)
	fmt.Printf("dimensions %#v, startPos %v\n", dimensions, startPos)
	digged := make([][]bool, dimensions.row)
	for i := range digged {
		digged[i] = make([]bool, dimensions.col)
	}
	currentPosition := startPos
	for _, instruction := range instructions {
		for i := uint64(0); i < instruction.amount; i++ {
			digged[currentPosition.row][currentPosition.col] = true
			currentPosition = currentPosition.add(instruction.direction.toVector())
			// currently ignoring color
		}
	}
	return digged
}

// horizontal ray/intersection or widen cols/rows and flood fill?

func fill(holes [][]bool) [][]bool {
	for i, row := range holes[:len(holes)-1] { // skip last row, cannot fill anything there
		if i == 0 {
			continue
		}

		// previousDigged := false
		// below := false // toggle everytime below is digged
		isIn := false
		for j, digged := range row[:len(row)-1] {
			belowDigged := holes[i+1][j]
			// if belowDigged {
			// 	below = !below
			// }
			// if j == 0 {
			// 	continue
			// }
			if digged && belowDigged {
				isIn = !isIn
			}
			if isIn {
				holes[i][j] = true
			}

			// previousDigged = digged
		}
	}
	return holes
}

func countTrue(mask [][]bool) (total int) {
	for _, row := range mask {
		for _, v := range row {
			if v {
				total++
			}
		}
	}
	return total
}

func part1(instructions []Instruction) (total int) {
	edge := dig(instructions)
	// pprint(edge)
	// println()
	fmt.Println("caclulated edges")
	holes := fill(edge)
	fmt.Println("caclulated fill")
	// pprint(holes)
	return countTrue(holes)

}

func fixInstructions(instructions []Instruction) []Instruction {
	newInstructions := make([]Instruction, len(instructions))
	for i, inst := range instructions {
		hexAmount := inst.color[2:7]
		direction := inst.color[7]
		a, err := strconv.ParseUint(hexAmount, 16, 64)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("i %03d: hexAmount %q, direction %q\n", i, hexAmount, direction)
		d := right
		switch direction {
		case '0':
			d = right
		case '1':
			d = down
		case '2':
			d = left
		case '3':
			d = up
		default:
			panic("not recognized direction")
		}
		newInst := Instruction{
			amount:    a,
			direction: d,
		}
		newInstructions[i] = newInst
	}
	return newInstructions
}

func part2(instructions []Instruction) (maxTotal int) {
	actualInstructions := fixInstructions(instructions)
	fmt.Printf("fixed instructions")
	// for _, i := range actualInstructions {
	// 	fmt.Printf("%#v\n", i)
	// }
	// fmt.Printf()
	// how to cache shared ways?
	return part1(actualInstructions) // brute force, would need slice-logic to optimize?
}
