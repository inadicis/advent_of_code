package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type rangedMap struct {
	source_start      int
	length            int
	destination_start int
}

type BySource []rangedMap

func (a BySource) Len() int           { return len(a) }
func (a BySource) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySource) Less(i, j int) bool { return a[i].source_start < a[j].source_start }

func main() {
	result := solution5("input.txt")
	fmt.Println(result)
}

func solution5(filename string) int {
	content, err := os.ReadFile(filename)
	// use a scanner next time to not hold everything in memory
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	text_blocks := strings.Split(strings.Trim(string(content), " \r\n"), "\r\n\r\n") // separate blocks

	// get starting seeds
	_, seed_numbers, _ := strings.Cut(text_blocks[0], ":")
	seeds := strings.Split(strings.Trim(seed_numbers, " "), " ")

	allMaps := [][]rangedMap{}
	for _, block := range text_blocks[1:] {
		rangedMaps := []rangedMap{}
		for _, line := range strings.Split(block, "\r\n")[1:] { // skipping the name
			numbers := strings.Split(line, " ")
			source, err := strconv.Atoi(numbers[1])
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			destination, err := strconv.Atoi(numbers[0])
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			length, err := strconv.Atoi(numbers[2])
			if err != nil {
				log.Fatal(err)
				panic(err)
			}

			newMap := rangedMap{
				source_start:      source,
				destination_start: destination,
				length:            length,
			}
			rangedMaps = append(rangedMaps, newMap)
		}
		sort.Sort(BySource(rangedMaps))
		allMaps = append(allMaps, rangedMaps)
	}
	locations := []int{}
	for _, seed := range seeds {
		current, err := strconv.Atoi(seed)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		for _, rangedMaps := range allMaps {
			// linear search of the correct range. A binary search should be possible though
			for _, rangedMap := range rangedMaps {
				diff := current - rangedMap.source_start
				if diff < 0 {
					break // no map existing for this number (because sorted) (assuming no overlaps)
				}
				if diff < rangedMap.length {
					current = rangedMap.destination_start + diff
					break // map found, go to next step
				}
				// if there is a map, it is further down , so just continue
			}
		}
		locations = append(locations, current)
	}
	min_location := locations[0]
	for _, location := range locations[1:] {
		if location < min_location {
			min_location = location
		}
	}
	return min_location
	// return min(locations...)
	// TODO why does that not work? min wants fixed number of parameters and not variatic????
}

// Learning Go: TODO s
// - understand packages, modules, imports, both when only locally and when shared packages online
// - understand defer keyword
