package main

import (
	"bytes"
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
	// fmt.Printf(" getExit(%s, direction: %v\n)", currentPipe, entranceDirection)
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
	visited = make(map[Position]Tile)
	queue := list.New()
	for _, direction := range []Direction{up, right, down, left} {
		position, exists := getNextPosition(maze, startPosition, direction)
		if exists {
			queue.PushBack(Tile{position: position, entranceDirection: direction, distance: 1, origin: direction})
		}
	}
	for i := 0; ; i++ {
		front := queue.Front()
		if front == nil {
			panic("Only dead ends! No cycle contains given startPosition")
		}
		currentQueuedTile := front.Value.(Tile)
		queue.Remove(front)

		tile, alreadyVisited := visited[currentQueuedTile.position]
		if alreadyVisited {
			direction1 := tile.origin
			direction2 := currentQueuedTile.origin
			startRune, found := pipeToRune[Pipe{direction1, direction2}]
			if !found {
				startRune, found = pipeToRune[Pipe{direction2, direction1}]
				if !found {
					fmt.Printf("direction1 %#v, direction2 %#v matches no rune", direction1, direction2)
					panic("No pipe found matching directions given")
				}

			}
			return currentQueuedTile, visited, startRune
		}

		// pipe, found := pipes[maze[currentQueuedPipe.position.row][currentQueuedPipe.position.col]]
		currentPipe := maze[currentQueuedTile.position.row][currentQueuedTile.position.col]
		exitDirection, found := getExitDirection(currentPipe, currentQueuedTile.entranceDirection)
		if found {
			nextPosition, exists := getNextPosition(maze, currentQueuedTile.position, exitDirection)
			if exists {
				queue.PushBack(Tile{position: nextPosition, entranceDirection: exitDirection, distance: currentQueuedTile.distance + 1, origin: currentQueuedTile.origin})
				visited[currentQueuedTile.position] = currentQueuedTile
			} else {
				fmt.Printf("step %d: Stopped branch ending with %#v as it pointed outside the maze", i, currentQueuedTile)
			}

		} else {
			fmt.Printf("step %d: Stopped branch ending with %#v as its position was not recognized as a pipe\n", i, currentQueuedTile)

		}

	}
}

func main() {
	maze, startPos, err := ExtractData("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	furthestTile, visitedTiles, startRune := BreadthFirstSearch(maze, startPos)
	fmt.Printf("Furtherst Position found at distance %d: %#v ", furthestTile.distance, furthestTile.position)
	// maze[startPos.row][startPos.col] = startRune
	visitedTiles[startPos] = Tile{position: startPos, distance: 0, origin: furthestTile.origin}
	area := countEnclosedTiles(maze, visitedTiles, startPos, startRune)
	fmt.Println(area)
}

func countEnclosedTiles(maze [][]rune, visitedTiles map[Position]Tile, startPos Position, startRune rune) (amount int) {
	// go through each tile in each row (horizontally). Check for each tile if it is "in" the circle.
	// for the current position, we can either be in the circle, at its edge (on a pipe), or outside.
	// We only get "in" the circle by "crossing" a pipe on the edge of the circle.
	// So we simply store if we are currently "in" the circle, and alternate this value everytime we "cross" the edge
	// and count how many tiles are in

	// to "enter" the area, each line has to "cross" a up-down | pipe.
	// this pipe could be "prolonged/divided" in multiple L--7 pipes or similar,
	// but the one extreme must be up, the other must be down to inverse In/Out (LJ would just be ignored)
	fmt.Println()
	var buffer bytes.Buffer
	maze[startPos.row][startPos.col] = startRune
	acceptedOrigins := runeToPipe[startRune]
	fmt.Println()
	for i, row := range maze {
		isIn := false
		crossedUp := false
		crossedDown := false
		for j, r := range row {
			tile, visited := visitedTiles[Position{i, j}]
			if visited && (tile.origin == acceptedOrigins.entrance || tile.origin == acceptedOrigins.exit) {

				buffer.WriteRune(r)
				switch {
				case r == '|':
					isIn = !isIn
				case r == 'J' || r == 'L':
					crossedUp = !crossedUp
				case r == 'F' || r == '7':
					crossedDown = !crossedDown
				}
				if crossedDown && crossedUp {
					isIn = !isIn
					crossedDown = false
					crossedUp = false
				}

			} else if isIn {
				amount++
				buffer.WriteRune('o')
			} else {
				buffer.WriteRune(' ')
			}

		}
		buffer.WriteString("\n")
	}
	visualization := buffer.String()
	// fmt.Print(visualization)
	for i, row := range strings.Split(visualization, "\n") {
		fmt.Printf("%03d:   %s\n", i, row)
	}
	return amount
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
