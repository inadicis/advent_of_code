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

func main() {
	text, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	mirrors := cleanupData(text)
	pprint(mirrors)
	result := part1(mirrors)
	fmt.Printf("\nFinal result: %d", result)

}

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
	if char == '.' { // TODO simplify if-elses
		outcome = append(outcome, d)
	} else {
		mirror := getMirror(char)
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
		}
	}
	return outcome

	// switch char {
	// case '-':
	// 	if !d.upDown {
	// 		outcome[0] = d
	// 	} else {
	// 		outcome[0] = Direction{}
	// 	}
	// }
}

type Tile struct {
	row       int
	col       int
	direction Direction
}

func part1(mirrors [][]rune) (total int) {
	visited := make(map[Tile]bool)
	activated := make([][]bool, len(mirrors))
	for i, row := range mirrors {
		activated[i] = make([]bool, len(row))
	}
	fmt.Printf("Len rows: %d, Len cols %d", len(mirrors), len(mirrors[0]))
	queue := list.New()
	queue.PushBack(Tile{0, 0, Direction{}})
	step := 0
	for e := queue.Front(); e != nil; e = e.Next() {
		tile := e.Value.(Tile)
		if visited[tile] {
			continue // cycle detected (same position and direction) ==> skip
		}
		activated[tile.row][tile.col] = true
		visited[tile] = true
		char := mirrors[tile.row][tile.col]
		fmt.Printf(" %02d: current char: %q, %#v\n", step, char, tile)
		nextDirections := reflections(tile.direction, char)
		fmt.Printf(" next: %#v\n", nextDirections)
		for _, d := range nextDirections {
			row := tile.row + d.vector().row
			col := tile.col + d.vector().col
			fmt.Printf(" | row %d, col %d...\n", row, col)
			if row < 0 || row >= len(mirrors) || col < 0 || col >= len(mirrors[0]) {
				fmt.Printf(" | skipping: len(rows) %d, len(cols) %d\n", len(mirrors), len(mirrors[0]))
				continue
			}
			queue.PushBack(Tile{row: row, col: col, direction: d})
		}
		step++
		// switch {
		// 	case char == ""
		// }
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

func part2(instructions []string) (total int) {
	return total
}
