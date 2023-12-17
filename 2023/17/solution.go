package main

import (
	"container/heap"
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
	bytes, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	maze := cleanupData(bytes)
	// pprint(maze)
	result := part1(maze)
	fmt.Printf("\nFinal result: %d", result)

}

// returns a 2D rune matrix (assumes unix-like newlines \n)
func cleanupData(text []byte) [][]int {
	lines := strings.Split(strings.TrimSpace(string(text)), "\n") // unix files
	rows := make([][]int, len(lines))
	for i, line := range lines {
		row := make([]int, len(line))
		for j, char := range line {
			row[j] = int(char - '0')
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

func (v Vector) add(v2 Vector) Vector {
	return Vector{v.row + v2.row, v.col + v2.col}
}

func (v Vector) isOutOfBounds(maze [][]int) bool {
	return v.row < 0 || v.col < 0 || v.row >= len(maze) || v.col >= len(maze[v.row])
}

// Combination of position and direction
// type State struct {
// 	position  Vector
// 	direction ManhattanDirection
// 	streak    int // amount of steps in this direction
// }

func possibleNextStates(state State) []State {
	newStates := make([]State, 0, 3)
	if state.streak < 3 {
		s := State{
			position:  state.position.add(state.direction.toVector()),
			direction: state.direction,
			streak:    state.streak + 1}
		newStates = append(newStates, s)
	}
	// left and right
	d1 := ManhattanDirection{isVertical: !state.direction.isVertical, sign: true}
	d2 := ManhattanDirection{isVertical: !state.direction.isVertical, sign: false}
	newStates = append(newStates,
		State{
			position:  state.position.add(d1.toVector()),
			direction: d1,
			streak:    1,
		})
	newStates = append(newStates,
		State{
			position:  state.position.add(d2.toVector()),
			direction: d1,
			streak:    1,
		})
	return newStates
}

type State struct {
	position Vector
	distance int
	// here priority = inverse of distance
	index     int
	direction ManhattanDirection
	streak    int
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*State)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *State, distance int) {
	item.distance = distance
	heap.Fix(pq, item.index)
}

// func (position Vector) toInitialQueueIndex(maze [][]int) int {
// 	return position.row*len(maze[0]) + position.col
// }

func findShortestPath(maze [][]int, source State) (total int) {
	// djikstra from source
	// using a simple priority queue -> TODO implement fibonacci heap
	infinite := 1
	distances := make([][]int, len(maze))
	visited := make([][]bool, len(maze))
	for _, row := range maze {
		for _, weight := range row {
			infinite += weight
		}
	}
	// pq := make(PriorityQueue, len(maze)*len(maze[0]))
	for i, row := range maze {
		visited[i] = make([]bool, len(row))
		distances[i] = make([]int, len(row))
		for j := range distances {
			if i != source.position.row || j != source.position.col {
				distances[i][j] = infinite
			}

			// queueIndex := i*len(row) + j
			// pq[queueIndex] = &State{
			// 	distance: infinite,
			// 	index:    queueIndex,
			// 	position: Vector{i, j},
			// }
		}
	}
	pq := PriorityQueue{&source}

	heap.Init(&pq)

	// queue := list.New()
	// queue.PushFront(currentState)
	// for e := queue.Front(); e != nil; e = queue.Front() {

	// currentDirection := ManhattanDirection{}
	// currentStreak := 0
	fmt.Printf("Initialized: %#v\n %#v\n", source, distances)
	for pq.Len() > 0 {
		// state := e.Value.(State)
		state := *heap.Pop(&pq).(*State)
		fmt.Printf("| Working with State %#v\n", state)
		if visited[state.position.col][state.position.col] {

			fmt.Print("- Skipped because visited already\n")
			continue
		}
		visited[state.position.col][state.position.col] = true
		currentDistance := distances[state.position.col][state.position.col]
		fmt.Printf("| Current Distance %d\n", currentDistance)
		for _, newState := range possibleNextStates(state) {

			fmt.Printf("|| Adjacent state %#v\n", newState)
			row := newState.position.row
			col := newState.position.col
			if newState.position.isOutOfBounds(maze) || visited[row][col] {
				fmt.Print("|- Skipped because OOB\n")
				continue
			}
			previousDist := distances[row][col]
			newState.distance = currentDistance + maze[row][row]
			if newState.distance < previousDist {
				fmt.Printf("|| new distance %d for adjacent is tinier than previous %d\n", newState.distance, previousDist)
				distances[row][col] = newState.distance
			}
			pq.Push(&newState)

		}
		// careful possible directions, not out of bounds, not visited
	}

	return total
}
func part1(maze [][]int) (total int) {
	return findShortestPath(maze, State{
		position:  Vector{0, 0},
		distance:  0,
		direction: right,
		streak:    0,
	})
}

func part2(maze [][]rune) (maxTotal int) {
	// how to cache shared ways?

	return maxTotal
}
