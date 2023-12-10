package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strings"
)

type Position struct {
	row int
	col int
}

type Direction Position

var (
	up    = Direction{-1, 0}
	down  = Direction{1, 0}
	right = Direction{0, 1}
	left  = Direction{0, -1}
)

type Pipe struct {
	char     rune
	entrance Direction
	exit     Direction
}

var pipes map[rune]Pipe = map[rune]Pipe{
	'|': {'|', up, down},
	'-': {'-', left, right},
	'L': {'L', up, right},
	'J': {'J', up, left},
	'7': {'7', left, down},
	'F': {'F', right, down},
}

func getExitDirection(currentPipe rune, entranceDirection Direction) (exit Direction, found bool) {
	pipe, found := pipes[currentPipe]
	if !found {
		return Direction{}, found
	}
	if entranceDirection == pipe.entrance {
		return pipe.exit, true
	} else if entranceDirection == pipe.exit {
		return pipe.entrance, true
	} else {
		panic("pipe did not match current direction")
	}
}

func getNextPosition(maze [][]rune, position Position, direction Direction) (nextPosition Position, exists bool) {
	row := position.row + direction.row
	col := position.col + direction.col
	if row < 0 || col < 0 || row >= len(maze) || col >= len(maze[0]) {
		return Position{}, false
	}
	return Position{row, col}, true
}

type QueuedPipe struct {
	position          Position
	entranceDirection Direction
	distance          int
}

func BreadthFirstSearch(maze [][]rune, startPosition Position) (furthestPositionInCycle Position, distance int) {
	visited := make(map[Position]int)
	// queue := make([]Position, 4)
	queue := list.New()
	for _, direction := range []Direction{up, left, right, down} {
		position, exists := getNextPosition(maze, startPosition, direction)
		if exists {
			queue.PushBack(QueuedPipe{position: position, entranceDirection: direction, distance: 1})
		}
		// queue[i] = Position{row, col}
	}
	for i := 0; ; i++ {
		fmt.Printf("Investigating pipe %d", i)
		front := queue.Front()
		if front == nil {
			panic("Only dead ends!")
		}
		currentQueuedPipe := front.Value.(QueuedPipe)
		queue.Remove(front)

		distance, alreadyVisited := visited[currentQueuedPipe.position]
		if alreadyVisited {
			return currentQueuedPipe.position, distance
		}

		// pipe, found := pipes[maze[currentQueuedPipe.position.row][currentQueuedPipe.position.col]]
		currentPipe := maze[currentQueuedPipe.position.row][currentQueuedPipe.position.col]
		exitDirection, found := getExitDirection(currentPipe, currentQueuedPipe.entranceDirection)
		if found {
			nextPosition, exists := getNextPosition(maze, currentQueuedPipe.position, exitDirection)
			if exists {
				queue.PushBack(QueuedPipe{position: nextPosition, entranceDirection: exitDirection, distance: currentQueuedPipe.distance + 1})
				visited[currentQueuedPipe.position] = currentQueuedPipe.distance
			} else {
				fmt.Printf("step %d: Stopped branch ending with %#v as it pointed outside the maze", i, currentQueuedPipe)
			}

		} else {
			fmt.Printf("step %d: Stopped branch ending with %#v as its position was not recognized as a pipe", i, currentQueuedPipe)

		}

	}
}

// func findStartPos(maze [][]rune)Position{
// 	for i, row := range maze {
// 		for j, char := range row{
// 			if char == 'S'{
// 				return
// 			}
// 		}
// 	}
// }

func main() {
	maze, startPos, err := ExtractData("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	furthestPos, distance := BreadthFirstSearch(maze, startPos)
	fmt.Printf("Furtherst Position found at distance %d: %#v ", distance, furthestPos)

}

func ExtractData(filename string) ([][]rune, Position, error) {
	block, err := os.ReadFile(filename)
	if err != nil {
		return [][]rune{}, Position{}, err
	}
	var startPos Position
	lines := strings.Split(strings.TrimSpace(string(block)), "\r\n")
	rows := make([][]rune, len(lines))
	for i, line := range lines {
		row := make([]rune, len(line))
		for j, char := range line {
			row[j] = rune(char)
			if row[j] == 'S' {
				startPos = Position{row: i, col: j}
			}
		}
		rows[i] = row
	}
	return rows, startPos, nil
}

func day10(filename string) {

}
