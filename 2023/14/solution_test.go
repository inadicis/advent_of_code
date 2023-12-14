package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParts(t *testing.T) {
	tests := []struct {
		expected int
		filename string
		amount   int
		fn       func(string) int
	}{
		{
			expected: 136,
			filename: "test.txt",
			fn:       part1,
		},
		{
			expected: 64,
			filename: "test.txt",
			fn:       part2,
		},
	}
	for _, test := range tests {
		b, err := os.ReadFile(test.filename)
		assert.NoError(t, err, test.filename)
		assert.Equal(t, test.expected, test.fn(string(b)))
	}
}

func TestInternalFunc(t *testing.T) {
	tests := []struct {
		expected     []string
		input        []string
		amountCycles int
		fn           func([][]rune)
	}{
		{
			expected: []string{
				".....#....",
				"....#...O#",
				"...OO##...",
				".OO#......",
				".....OOO#.",
				".O#...O#.#",
				"....O#....",
				"......OOOO",
				"#...O###..",
				"#..OO#....",
			},
			input: []string{
				"O....#....",
				"O.OO#....#",
				".....##...",
				"OO.#O....O",
				".O.....O#.",
				"O.#..O.#.#",
				"..O..#O..O",
				".......O..",
				"#....###..",
				"#OO..#....",
			},
			amountCycles: 1,
			fn:           spinRocks,
		},
		{
			expected: []string{
				".....#....",
				"....#...O#",
				".....##...",
				"..O#......",
				".....OOO#.",
				".O#...O#.#",
				"....O#...O",
				".......OOO",
				"#..OO###..",
				"#.OOO#...O",
			},
			input: []string{
				"O....#....",
				"O.OO#....#",
				".....##...",
				"OO.#O....O",
				".O.....O#.",
				"O.#..O.#.#",
				"..O..#O..O",
				".......O..",
				"#....###..",
				"#OO..#....",
			},
			amountCycles: 2,
			fn:           spinRocks,
		},
		{
			expected: []string{
				".....#....",
				"....#...O#",
				".....##...",
				"..O#......",
				".....OOO#.",
				".O#...O#.#",
				"....O#...O",
				".......OOO",
				"#...O###.O",
				"#.OOO#...O",
			},
			input: []string{
				"O....#....",
				"O.OO#....#",
				".....##...",
				"OO.#O....O",
				".O.....O#.",
				"O.#..O.#.#",
				"..O..#O..O",
				".......O..",
				"#....###..",
				"#OO..#....",
			},
			amountCycles: 3,
			fn:           spinRocks,
		},
	}
	for _, test := range tests {
		result := StringsToRunes(test.input)
		for i := 0; i < test.amountCycles; i++ {
			spinRocks(result)
		}
		assert.Equal(t, StringsToRunes(test.expected), result)
	}
}

func TestCountWeight(t *testing.T) {
	tests := []struct {
		expected int
		input    []string
	}{
		{
			expected: 136,
			input: []string{
				"OOOO.#.O..",
				"OO..#....#",
				"OO..O##..O",
				"O..#.OO...",
				"........#.",
				"..#....#.#",
				"..O..#.O.O",
				"..O.......",
				"#....###..",
				"#....#....",
			},
		},
	}
	for _, test := range tests {
		result := countWeight(StringsToRunes(test.input))
		assert.Equal(t, test.expected, result)
	}

}
