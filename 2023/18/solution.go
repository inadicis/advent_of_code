package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"slices"
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
	result2 := part1(instructions)
	fmt.Printf("Actual result: %d\n", result2)
	result := part1Optimized(instructions)
	fmt.Printf("\nFinal result: %d\n", result)
}

type Instruction struct {
	direction ManhattanDirection
	amount    int64
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
			amount:    int64(a),
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
	n := int64(1)
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
	row int64 // (x) 0, -1 or 1 for directions
	col int64 // (-y)  0, -1 or 1 for directions
}

func (v Vector) add(v2 Vector) Vector {
	return Vector{v.row + v2.row, v.col + v2.col}
}

func (v Vector) isOutOfBounds(maze [][]int) bool {
	return v.row < 0 || v.col < 0 || v.row >= int64(len(maze)) || v.col >= int64(len(maze[v.row]))
}

func getInitialState(instructions []Instruction) (dimensions Vector, startPosition Vector) {
	currentPos := Vector{}

	var minRow, minCol, maxRow, maxCol int64
	for _, instruction := range instructions {
		dirVector := instruction.direction.toVector()
		currentPos = currentPos.add(Vector{dirVector.row * instruction.amount, dirVector.col * instruction.amount})
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
		for i := int64(0); i < instruction.amount; i++ {
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
	pprint(holes)
	return countTrue(holes)

}

type VerticalEdge struct {
	col      int64
	rowStart int64
	rowEnd   int64
}

func getVerticalEdges(instructions []Instruction) ([]VerticalEdge, Vector) {
	e := make([]VerticalEdge, 0, len(instructions)/2+1)
	dimensions, startPos := getInitialState(instructions)
	// fmt.Printf("dimensions %#v, startPos %v\n", dimensions, startPos)
	currentPos := startPos
	for _, instruction := range instructions {
		col := currentPos.col
		startRow := currentPos.row
		dirVector := instruction.direction.toVector()
		currentPos = currentPos.add(Vector{dirVector.row * instruction.amount, dirVector.col * instruction.amount})
		endRow := currentPos.row
		if instruction.direction.isVertical {
			if endRow < startRow {
				startRow, endRow = endRow, startRow
			}
			e = append(e, VerticalEdge{col, startRow, endRow})
		}
	}

	return e, dimensions
}

type State struct {
	edge        VerticalEdge
	changeIndex int64
	isStart     bool
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { // lower value has higher priority
	return pq[i].changeIndex < pq[j].changeIndex
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	// pq[i].index = i
	// pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	// n := len(*pq)
	// item :=
	// item.index = n
	*pq = append(*pq, x.(*State))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	// old[n-1] = nil  // avoid memory leak
	// item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func splitEdgesIntoSharedVerticalBatches(edges []VerticalEdge) [][]VerticalEdge {
	// slices.SortFunc(edges, func(a VerticalEdge, b VerticalEdge) int { return int(a.rowStart - b.rowStart) })
	pq := make(PriorityQueue, 2*len(edges))
	for i, e := range edges {
		pq[2*i] = &State{e, e.rowStart, true}
		pq[2*i+1] = &State{e, e.rowEnd, false}
	}
	heap.Init(&pq)
	// fmt.Printf("Priority queue initialized, len: %d\n", len(pq))
	batches := make([][]VerticalEdge, 0, len(edges)) // at least one batch per vertical edge if no intersections
	previousEdges := make([]VerticalEdge, 0)
	// currentIndex := 0
	for pq.Len() > 0 {
		state := heap.Pop(&pq).(*State)
		newBatch := make([]VerticalEdge, 0, len(previousEdges))
		updatedPreviousEdges := make([]VerticalEdge, 0, len(previousEdges))
		fmt.Printf("New change at row %d, isStart %t, len(previousEdges)=%d\n", state.changeIndex, state.isStart, len(previousEdges))
		fmt.Printf("%v\n", state.edge)
		for _, previousEdge := range previousEdges {
			// because of heap condition, they all go at least as far as state.changeIndex
			if state.changeIndex > previousEdge.rowStart {
				newBatch = append(newBatch, VerticalEdge{col: previousEdge.col, rowStart: previousEdge.rowStart, rowEnd: state.changeIndex})
			}
			if previousEdge.rowEnd > state.changeIndex {
				updatedPreviousEdges = append(updatedPreviousEdges, VerticalEdge{
					col:      previousEdge.col,
					rowStart: state.changeIndex,
					rowEnd:   previousEdge.rowEnd,
				})
			}
			previousEdges = updatedPreviousEdges
		}
		if state.isStart {
			fmt.Printf("Appending edge %d->%d to previousEdges\n", state.edge.rowStart, state.edge.rowEnd)
			previousEdges = append(previousEdges, state.edge)
		}
		fmt.Printf("Previous Edges: %v\n", previousEdges)
		fmt.Printf("New Batch: %v\n", newBatch)
		if len(newBatch) > 0 {
			batches = append(batches, newBatch)
		}
		// TODO careful about edges that change at the same row
		// for _, newEdge := range edges {
		// for _, previousEdge := range previousEdges{
		// 	if previousEdge.rowEnd <
		// }
	}
	return batches
}

func part1Optimized(instructions []Instruction) (total int64) {
	verticalEdges, dimensions := getVerticalEdges(instructions)
	// fmt.Printf("%#v\n", verticalEdges)
	fmt.Printf("%#v\n", dimensions)
	// for i, e := range verticalEdges {
	// fmt.Printf("i %03d: col %d, row %d -> %d\n", i, e.col, e.rowStart, e.rowEnd)
	// }
	batches := splitEdgesIntoSharedVerticalBatches(verticalEdges)
	var previousBatch []VerticalEdge
	for i, b := range batches {
		slices.SortFunc(b, func(a, b VerticalEdge) int { return int(a.col - b.col) })
		fmt.Printf("Batch i %03d: (%d -> %d) : [ ", i, b[0].rowStart, b[0].rowEnd)
		for _, e := range b {
			fmt.Printf("col %d,", e.col)
		}
		fmt.Println(" ]")
		for j := 0; j < len(b); j += 2 {
			total += (b[j+1].col - b[j].col + 1) * (b[j].rowEnd - b[j].rowStart + 1)
			// counts twice the intersection between the rectangles one above the other
		}
		if i != 0 {
			fmt.Printf(" %03d: Eliminating overlaps with previous batch: ", i)
			for j, k := 0, 0; j < len(b) && k < len(previousBatch); j += 2 {
				// eliminate overlaps with previous batches
				// as there is no horizonzal overlap, we only need to iterate by pairing segments from
				// previous batch and current batch
				// (overlap is only one row high, so we do simple segment arithmetic here)
				pStart := previousBatch[k].col
				pEnd := previousBatch[k+1].col
				cStart := b[j].col
				cEnd := b[j+1].col
				duplicated := max(0, min(cEnd, pEnd)-max(cStart, pStart)+1)
				total -= duplicated
				if pEnd < cEnd {
					k += 2
					j -= 2
				}
				fmt.Print(duplicated, " ")
			}
			fmt.Println()
		}
		// var batchIndexPrevious, batchIndexCurrent int
		// for {
		// 	previous := previousBatch[batchIndexPrevious]
		// 	current := b[batchIndexCurrent]
		// 	if
		// }

		previousBatch = b
	}

	// pprint(edge)
	// println()
	// fmt.Println("caclulated edges")
	// holes := fill(edge)
	// fmt.Println("caclulated fill")
	// pprint(holes)
	// return countTrue(holes)
	return total

}

func fixInstructions(instructions []Instruction) []Instruction {
	newInstructions := make([]Instruction, len(instructions))
	for i, inst := range instructions {
		hexAmount := inst.color[2:7]
		direction := inst.color[7]
		a, err := strconv.ParseInt(hexAmount, 16, 64)
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

func part2(instructions []Instruction) (maxTotal int64) {
	actualInstructions := fixInstructions(instructions)
	fmt.Printf("fixed instructions")
	return part1Optimized(actualInstructions)
}
