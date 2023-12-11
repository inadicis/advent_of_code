package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	universe, galaxies, err := ExtractData("test_data.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(universe)
	fmt.Println(galaxies)
	actualUnivers, actualGalaxies := ExtandGaps(universe, galaxies)
	// fmt.Println(actualUnivers)
	fmt.Println(actualGalaxies)
	for _, row := range actualUnivers {
		// fmt.Printf("%03d: %v\n", i, row)
		for _, char := range row {
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
	// distances := CalculateDistances(data)
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

func ExtandGaps(universe [][]rune, galaxies [][]int) ([][]rune, [][]int) {
	rowHasGalaxy := make([]bool, len(universe))
	colHasGalaxy := make([]bool, len(universe[0]))
	for _, galaxy := range galaxies {
		// fmt.Printf("universe size: %dx%d. current galaxy %v", len(universe), len(universe[0]), galaxy)
		rowHasGalaxy[galaxy[0]] = true
		colHasGalaxy[galaxy[1]] = true
	}
	fmt.Printf("has row galaxy: %v\n", rowHasGalaxy)
	fmt.Printf("has col galaxy: %v\n", colHasGalaxy)
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
	actualRowIndex := 0
	actualUniverse := make([][]rune, len(universe)+addedRows)
	actualGalaxies := make([][]int, len(galaxies))
	galaxyCount := 0
	for i, row := range universe {
		if !rowHasGalaxy[i] {
			fmt.Printf(" Adding a row at i %d\n", i)
			actualUniverse[actualRowIndex] = make([]rune, len(universe[0])+addedCols)
			for i := range actualUniverse[actualRowIndex] {
				actualUniverse[actualRowIndex][i] = ' '
			}
			actualRowIndex++
		}
		actualColIndex := 0
		newRow := make([]rune, len(universe[0])+addedCols)
		for j, char := range row {
			if !colHasGalaxy[j] {

				fmt.Printf("  Adding a col at j %d (actual index: %d)\n", j, actualColIndex)
				newRow[actualColIndex] = ' '
				actualColIndex++
			}
			newRow[actualColIndex] = char
			actualColIndex++
			if char == '#' {
				actualGalaxies[galaxyCount] = []int{actualRowIndex, actualColIndex}
				galaxyCount++
			}
		}
		actualUniverse[actualRowIndex] = newRow
		actualRowIndex++
	}
	return actualUniverse, actualGalaxies
}
