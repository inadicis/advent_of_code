package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strings"
)

func pprint(lines [][]rune) {
	for _, row := range lines {
		for _, char := range row {
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
}

var ( // GOLANG? how to make this const? const structs possible?
	up    = Direction{true, true}
	down  = Direction{true, false}
	right = Direction{false, false}
	left  = Direction{false, true}
)

func main() {
	text, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	maze := cleanupData(text)
	pprint(maze)
	result := part2(maze)
	fmt.Printf("\nFinal result: %d", result)

}

// GOLANG? -> how to split by newline independently of if it is on mac \n or windows \r\n?
func cleanupData(text []byte) [][]rune {
	lines := strings.Split(strings.TrimSpace(string(text)), "\n")
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

/*


 */

type Direction struct {
	upDown  bool // Up/Down or Right/Left
	reverse bool // false: Down or Right, true: Up or Left
}

func (d Direction) vector() (v DirectionVector) {
	n := 1
	if d.reverse {
		n = -1
	}
	if d.upDown {
		v.row = n
	} else {
		v.col = n
	}
	return v
}

type DirectionVector struct {
	row int // 0, -1 or 1
	col int
}

func (v DirectionVector) direction() (d Direction) {
	// assumes direction vector has exactly one of row/col != 0
	if v.col != 0 {
		d.upDown = true
	}

	if v.row < 0 || v.col < 0 {
		d.reverse = true
	}
	return d
}

type Mirror struct {
	straight bool
	inverse  bool
}

var mirrors map[rune]Mirror = map[rune]Mirror{
	'-':  {straight: true, inverse: false},
	'|':  {straight: true, inverse: true},
	'/':  {straight: false, inverse: true},
	'\\': {straight: false, inverse: false},
}

func getMirror(m rune) Mirror {
	switch m {
	case '-':
		return Mirror{straight: true, inverse: false}
	case '|':
		return Mirror{straight: true, inverse: true}
	case '/':
		return Mirror{straight: false, inverse: true}
	case '\\':
		return Mirror{straight: false, inverse: false}
	}
	panic("Did not recognize mirror")
}

func reflections(d Direction, char rune) []Direction {
	/*
		>|^v    >/^  >\v  >->

		^v|<    ^\<  v/<  <-<

		v   ^	v	v	>	<
		|	|	/	\	/	\
		v	^	<	>	^	^

		/ \ ==> upDown = !UpDown, / ==> reverse = !reverse
		-
	*/
	outcome := make([]Direction, 0, 2)

	mirror, isMirror := mirrors[char]
	if !isMirror {
		return append(outcome, d)
	}
	if mirror.straight {
		if mirror.inverse == d.upDown {
			outcome = append(outcome, d)
		} else {
			outcome = append(outcome, Direction{!d.upDown, true})
			outcome = append(outcome, Direction{!d.upDown, false})
		}
	} else {
		if mirror.inverse {
			outcome = append(outcome, Direction{!d.upDown, !d.reverse})
		} else {
			outcome = append(outcome, Direction{!d.upDown, d.reverse})
		}
		// outcome = append(outcome, Direction{!d.upDown, d.reverse != mirror.inverse})

	}

	return outcome
}

type Tile struct {
	row       int
	col       int
	direction Direction
}

func getAmountActivated(mirrors [][]rune, startTile Tile) (total int) {
	visited := make(map[Tile]bool)
	activated := make([][]bool, len(mirrors))
	for i, row := range mirrors {
		activated[i] = make([]bool, len(row))
	}
	queue := list.New()
	queue.PushBack(startTile)
	step := 0
	for e := queue.Front(); e != nil; e = e.Next() {
		tile := e.Value.(Tile)
		if visited[tile] {
			continue // cycle detected (same position and direction) ==> skip
		}
		activated[tile.row][tile.col] = true
		visited[tile] = true
		char := mirrors[tile.row][tile.col]
		nextDirections := reflections(tile.direction, char)
		for _, d := range nextDirections {
			row := tile.row + d.vector().row
			col := tile.col + d.vector().col
			if row < 0 || row >= len(mirrors) || col < 0 || col >= len(mirrors[0]) {
				continue
			}
			queue.PushBack(Tile{row: row, col: col, direction: d})
		}
		step++
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
	return getAmountActivated(maze, Tile{0, 0, Direction{}})
}

func part2(maze [][]rune) (maxTotal int) {
	for row := 0; row < len(maze); row++ {
		leftToRightTotal := getAmountActivated(maze, Tile{row, 0, right})
		rightToLeftTotal := getAmountActivated(maze, Tile{row, len(maze[row]) - 1, left})
		maxTotal = max(maxTotal, leftToRightTotal, rightToLeftTotal)
	}

	for col := 0; col < len(maze[0]); col++ {
		upToDownTotal := getAmountActivated(maze, Tile{0, col, up})
		downToUpTotal := getAmountActivated(maze, Tile{len(maze) - 1, col, down})
		maxTotal = max(maxTotal, upToDownTotal, downToUpTotal)
	}
	return maxTotal
}
