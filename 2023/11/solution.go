package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	universe, galaxies, err := ExtractData("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	_, rowGaps, colGaps := IdentifyGaps(universe, galaxies)
	// fmt.Println(actualUnivers)
	for _, row := range universe {
		// fmt.Printf("%03d: %v\n", i, row)
		for _, char := range row {
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
	gapFactor := 1000000
	distances := CalculateDistances(galaxies, rowGaps, colGaps, gapFactor)
	sum := 0
	for _, d := range distances {
		sum += d
	}
	fmt.Print(sum)
}

func ExtractData(filename string) ([][]rune, [][]int, error) {
	block, err := os.ReadFile(filename)
	if err != nil {
		return [][]rune{}, [][]int{}, err
	}
	lines := strings.Split(strings.TrimSpace(string(block)), "\r\n")
	rows := make([][]rune, len(lines))
	galaxyIndexes := [][]int{}
	for i, line := range lines {
		row := make([]rune, len(line))
		for j, char := range line {
			row[j] = rune(char)
			if row[j] == '#' {
				galaxyIndexes = append(galaxyIndexes, []int{i, j})
			}
		}
		rows[i] = row
	}
	return rows, galaxyIndexes, nil
}

func IdentifyGaps(universe [][]rune, galaxies [][]int) ([][]rune, []bool, []bool) {
	// Replaces all gaps with a space (instead of point)
	rowHasGalaxy := make([]bool, len(universe))
	colHasGalaxy := make([]bool, len(universe[0]))
	for _, galaxy := range galaxies {
		// fmt.Printf("universe size: %dx%d. current galaxy %v", len(universe), len(universe[0]), galaxy)
		rowHasGalaxy[galaxy[0]] = true
		colHasGalaxy[galaxy[1]] = true
	}
	// fmt.Printf("has row galaxy: %v\n", rowHasGalaxy)
	// fmt.Printf("has col galaxy: %v\n", colHasGalaxy)
	addedRows := 0
	for _, hasGalaxy := range rowHasGalaxy {
		if !hasGalaxy {
			addedRows++
		}
	}

	addedCols := 0
	for _, hasGalaxy := range colHasGalaxy {
		if !hasGalaxy {
			addedCols++
		}
	}
	for i, row := range universe {
		if !rowHasGalaxy[i] {
			// fmt.Printf(" Adding a row at i %d\n", i)
			for j := range row {
				universe[i][j] = ' '
			}
		} else {
			for j := range row {
				if !colHasGalaxy[j] {

					// fmt.Printf("  Adding a col at j %d \n", j)
					universe[i][j] = ' '
				}
			}
		}

	}
	return universe, rowHasGalaxy, colHasGalaxy
}

func CalculateDistances(galaxies [][]int, rowGaps []bool, colGaps []bool, gapFactor int) []int {
	distances := []int{}
	for i, galaxy := range galaxies {
		for j, otherGalaxy := range galaxies {
			if j >= i {
				continue
			}
			// manhattan distances: d = deltaX + deltaY
			minRow := galaxy[0]
			minCol := galaxy[1]
			maxRow := otherGalaxy[0]
			maxCol := otherGalaxy[1]
			if maxRow < minRow {
				minRow, maxRow = maxRow, minRow
			}
			if maxCol < minCol {
				minCol, maxCol = maxCol, minCol
			}
			deltaRow := maxRow - minRow
			deltaCol := maxCol - minCol

			gapsRow := rowGaps[minRow:maxRow]
			gapsCol := colGaps[minCol:maxCol]
			// we could exclude both rows were the galaxy are, but need to check if delta > 2
			amountGaps := 0

			for _, hasGalaxy := range gapsRow {
				if !hasGalaxy {
					amountGaps++
				}
			}

			for _, hasGalaxy := range gapsCol {
				if !hasGalaxy {
					amountGaps++
				}
			}
			actualDistance := deltaRow + deltaCol + amountGaps*(gapFactor-1) // gaps were already counted once in deltas
			distances = append(distances, actualDistance)
		}
	}
	return distances
}
