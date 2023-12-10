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
	// char     rune
	entrance Direction
	exit     Direction
}

var runeToPipe map[rune]Pipe = map[rune]Pipe{
	// always from in order which one comes first, which one second: up > right > down > left
	'|': {up, down},
	'L': {up, right},
	'J': {up, left},
	'F': {right, down},
	'-': {right, left},
	'7': {down, left},
}

var pipeToRune map[Pipe]rune = map[Pipe]rune{
	{up, down}:    '|',
	{up, right}:   'L',
	{up, left}:    'J',
	{right, down}: 'F',
	{right, left}: '-',
	{down, left}:  '7',
}

func getExitDirection(currentPipe rune, entranceDirection Direction) (exit Direction, found bool) {
	fmt.Printf(" getExit(%s, direction: %v)", currentPipe, entranceDirection)
	reversedDirection := Direction{-entranceDirection.row, -entranceDirection.col}
	pipe, found := runeToPipe[currentPipe]
	if !found {
		return Direction{}, false
	}
	if reversedDirection == pipe.entrance {
		return pipe.exit, true
	} else if reversedDirection == pipe.exit {
		return pipe.entrance, true
	} else {
		// panic("pipe did not match current direction")
		return Direction{}, false
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

type Tile struct {
	position          Position
	entranceDirection Direction
	distance          int
	origin            Direction // which adjacent tile of startPos was the origin of this traversal
}

// type VisitedTile struct {
// 	distance int
// 	origin   Direction // which adjacent tile of startPos was the origin of this traversal
// }

func BreadthFirstSearch(maze [][]rune, startPosition Position) (furthestTile Tile, visited map[Position]Tile, startPipe rune) {
	// could have just iterated over one way and come back to start
	// (as there are only 2 valid adjacent tiles to start this must work)
	// visited := make(map[Position]int)
	// queue := make([]Position, 4)
	queue := list.New()
	for _, direction := range []Direction{up, right, down, left} {
		position, exists := getNextPosition(maze, startPosition, direction)
		if exists {
			queue.PushBack(Tile{position: position, entranceDirection: direction, distance: 1, origin: direction})
		}
		// queue[i] = Position{row, col}
	}
	// fmt.Printf("Initialized BFS with queue: %#v\n", queue)
	for i := 0; ; i++ {
		// fmt.Printf(" Investigating pipe %d\n", i)
		front := queue.Front()
		if front == nil {
			panic("Only dead ends! No cycle contains given startPosition")
		}
		currentQueuedTile := front.Value.(Tile)
		queue.Remove(front)
		// fmt.Printf("  looking at pipe %#v\n", currentQueuedTile)

		tile, alreadyVisited := visited[currentQueuedTile.position]
		if alreadyVisited {
			direction1 := tile.origin
			direction2 := currentQueuedTile.origin
			startRune, found := pipeToRune[Pipe{direction1, direction2}]
			if !found {
				panic("No pipe found matching directions given")
			}
			return currentQueuedTile, visited, startRune
		}

		// pipe, found := pipes[maze[currentQueuedPipe.position.row][currentQueuedPipe.position.col]]
		currentPipe := maze[currentQueuedTile.position.row][currentQueuedTile.position.col]
		exitDirection, found := getExitDirection(currentPipe, currentQueuedTile.entranceDirection)
		if found {
			nextPosition, exists := getNextPosition(maze, currentQueuedTile.position, exitDirection)
			if exists {
				queue.PushBack(Tile{position: nextPosition, entranceDirection: exitDirection, distance: currentQueuedTile.distance + 1})
				visited[currentQueuedTile.position] = currentQueuedTile
			} else {
				fmt.Printf("step %d: Stopped branch ending with %#v as it pointed outside the maze", i, currentQueuedTile)
			}

		} else {
			fmt.Printf("step %d: Stopped branch ending with %#v as its position was not recognized as a pipe", i, currentQueuedTile)

		}

	}
}

func getBorders(maze [][]rune, startPosition Position) [][]rune {
	// initialize maze with only empty string tiles,
	// then fill only the tiles that are part of the cycle
	cleanMaze := make([][]rune, len(maze))
	for i, row := range maze {
		cleanMaze[i] = make([]rune, len(row))
	}
	currentPosition := startPosition
	for {

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
	// fmt.Print(getExitDirection('L', down))
	furthestTile, visitedTiles, startRune := BreadthFirstSearch(maze, startPos)
	fmt.Printf("Furtherst Position found at distance %d: %#v ", furthestTile.distance, furthestTile.position)

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
				fmt.Printf("Found Start position %#v\n", startPos)
			}
		}
		fmt.Println(row)
		rows[i] = row
	}
	return rows, startPos, nil
}
