package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
)

func pprint(lines [][]int) {
	for _, row := range lines {
		for _, char := range row {
			fmt.Printf("%03d ", char)
		}
		fmt.Println()
	}
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	// TODO find out why test2.txt has 12 at row 0, col 4 instead of 10 => typo d1/d2 in next possible directions
	// TODO find out how a way can be lower than the actual result?? (input.txt)
	// to debug input.txt i could change algorithm so that i store the final path to the last
	// cell then write a qa algo that checks if not more than 3 times in a row + sums are correct

	if err != nil {
		log.Fatal(err)
	}
	maze := cleanupData(bytes)
	pprint(maze)
	result := part1(maze)
	fmt.Printf("\nFinal result: %d", result)
	// 1294 too high
	// 1255 too low
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
			direction: d2,
			streak:    1,
		})
	// fmt.Printf("\\\\Possible next states: %#v\n", newStates)
	return newStates
}

type State struct {
	position Vector
	distance int
	// here priority = inverse of distance
	// index     int
	direction ManhattanDirection
	streak    int
}

type Visited struct {
	position  Vector
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

// update modifies the priority and value of an Item in the queue.
// func (pq *PriorityQueue) update(item *State, distance int) {
// 	item.distance = distance
// 	heap.Fix(pq, item.index)
// }

// func (position Vector) toInitialQueueIndex(maze [][]int) int {
// 	return position.row*len(maze[0]) + position.col
// }

func findShortestPath(maze [][]int, source State) (total int) {
	// djikstra from source
	// using a simple priority queue -> TODO implement fibonacci heap
	infinite := 1
	distances := make([][]int, len(maze))
	visited := make(map[Visited]bool)
	for _, row := range maze {
		for _, weight := range row {
			infinite += weight
		}
	}

	// pq := make(PriorityQueue, len(maze)*len(maze[0]))
	for i, row := range maze {
		distances[i] = make([]int, len(row))
		for j := range row {
			if i != source.position.row || j != source.position.col {
				distances[i][j] = infinite
			} else {
				distances[i][j] = source.distance
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
	// fmt.Printf("Initialized: %#v\n %#v\n", source, distances)
	for pq.Len() > 0 {
		// state := e.Value.(State)
		state := *heap.Pop(&pq).(*State)
		// fmt.Printf("| Working with State %#v\n", state)
		if visited[Visited{state.position, state.direction, state.streak}] {
			// fmt.Print("\\ Skipped because visited already with this exact direction/streak\n")
			continue
		}
		visited[Visited{state.position, state.direction, state.streak}] = true
		// currentDistance := distances[state.position.col][state.position.col]
		// fmt.Printf("| Current Distance %d\n", state.distance)
		for _, newState := range possibleNextStates(state) {

			// fmt.Printf("|| Adjacent state %#v\n", newState)
			row := newState.position.row
			col := newState.position.col
			debug := row == 0 && col == 4
			if debug {
				// fmt.Printf("\n previous state: %#v\n new state %#v\n", state, newState)
				// fmt.Println("added weight", maze[row][col])
			}
			if newState.position.isOutOfBounds(maze) || visited[Visited{newState.position, newState.direction, newState.streak}] {
				if debug {
					// fmt.Print("|\\ Skipped because OOB or visited\n")
				}
				continue
			}
			previousDist := distances[row][col]
			newD := state.distance + maze[row][col]
			if newD < previousDist {
				if debug {
					// fmt.Printf("|| new distance %d for adjacent is tinier than previous %d\n", newD, previousDist)
				}
				distances[row][col] = newD
			}
			s := State{
				newState.position,
				newD,
				newState.direction,
				newState.streak,
			}
			if debug {
				// fmt.Printf("|\\ Pushing new state into heap: %#v\n", s)
			}
			// heap.Push(&pq, &newState) // why is this not equivalent
			// fmt.Println(newState.position,
			// 	newState.distance,
			// 	newState.direction,
			// 	newState.streak)
			// fmt.Printf("heap before: %#v\n", pq)
			heap.Push(&pq, &s)
			// fmt.Printf("heap after: %#v\n", pq)
			// heap.Push(&pq, &newState)

		}
		// careful possible directions, not out of bounds, not visited
	}
	fmt.Println()
	pprint(distances)
	return distances[len(distances)-1][len(distances[0])-1]
}
func part1(maze [][]int) (total int) {
	// x := len(maze) -1
	// y := len(maze[0]) -1
	x, y := 0, 0

	return findShortestPath(maze, State{
		position: Vector{x, y},
		// distance:  maze[x][y],
		distance:  0, // already at this position, not entering it
		direction: down,
		streak:    0,
	})
}

func part2(maze [][]rune) (maxTotal int) {
	// how to cache shared ways?

	return maxTotal
}
