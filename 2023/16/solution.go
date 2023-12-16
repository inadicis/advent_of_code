package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strings"
)

func pprint(boxes []list.List) {

	for i, box := range boxes {
		s := ""
		for e := box.Front(); e != nil; e = e.Next() {
			s += fmt.Sprintf("[%v]", e.Value)
		}
		fmt.Printf(" box %03d: %s\n", i, s)
	}
}

func main() {
	text, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	instructions := cleanupData(text)
	result := part1(instructions)
	fmt.Printf("\nFinal result: %d", result)

}

func cleanupData(text []byte) [][]rune {
	lines := strings.Split(string(text), ",")
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
		d.upDown = True
	}

	if v.row < 0 || v.col < 0 {
		d.reverse = True
	}
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
	i         int
	j         int
	direction Direction
}

func part1(mirrors [][]rune) (total int) {
	activated := make([][]bool, len(mirrors))
	for i, row := range mirrors {
		activated[i] = make([]bool, len(row))
	}
	queue := list.New()
	queue.PushBack(Tile{0, 0, Direction{}})
	for e := queue.Front(); e != nil; e = e.Next() {
		tile := e.Value.(Tile)
		activated[tile.i][tile.j] = true
		char := mirrors[tile.i][tile.j]
		// switch {
		// 	case char == ""
		// }
	}
	return total
}

func part2(instructions []string) (total int) {
	return total
}
