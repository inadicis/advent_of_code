package main

import (
	"fmt"
	"log"
	"os"
	"slices"
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
	result, err := solution5("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

/*
Brute force solution is now implemented but will take minutes/hours to run
Optimization possibilities/ideas
- implement binary search instead of linear search for one map (minor/negligent impact)
- don't work each seed individually, find out ranges of seeds that have same destination (with offset)
  -> pass groups of seeds to next step. the list of groups will grow each step, but way less than actual amount of seeds (big impact)
- combine all maps into one to avoid repeated intermediate steps (moderate impact)

Delayed implementation to later/another day
*/

func solution5(filename string) (int, error) {
	content, err := os.ReadFile(filename)
	// use a scanner next time to not hold everything in memory
	if err != nil {
		return 0, err
	}
	text_blocks := strings.Split(strings.Trim(string(content), " \r\n"), "\r\n\r\n") // separate blocks

	// get starting seeds
	_, seed_numbers, _ := strings.Cut(text_blocks[0], ":")
	seeds := strings.Split(strings.Trim(seed_numbers, " "), " ")

	type seedRange struct {
		start  int
		length int
	}
	seedRanges := []seedRange{}
	previous := 0
	for i, n := range seeds {
		casted, err := strconv.Atoi(n)
		if err != nil {
			return 0, err
		}
		if i%2 == 0 {
			previous = casted
			continue
		}

		seedRanges = append(seedRanges, seedRange{start: previous, length: casted})
	}

	allMaps, err := buildMaps(text_blocks)

	if err != nil {
		return 0, err
	}

	locations := []int{}
	for _, seedRange := range seedRanges {
		fmt.Printf("seed range %+v\n", seedRange)
		for i := 0; i < seedRange.length; i++ {
			seed := seedRange.start + i
			// fmt.Printf(" seedrange %+v, seed: %d \n", seedRange, seed)
			current := seed
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

	}
	// min_location := locations[0]
	// for _, location := range locations[1:] {
	// 	if location < min_location {
	// 		min_location = location
	// 	}
	// }
	// return min_location
	return slices.Min(locations), nil
	// return min(locations...)
	// TODO why does that not work? min wants fixed number of parameters and not variatic????
}

func buildMaps(text_blocks []string) ([][]rangedMap, error) {

	allMaps := [][]rangedMap{}
	for _, block := range text_blocks[1:] {
		rangedMaps := []rangedMap{}
		for _, line := range strings.Split(block, "\r\n")[1:] { // skipping the name
			numbers := strings.Split(line, " ")
			source, err := strconv.Atoi(numbers[1])
			if err != nil {
				return [][]rangedMap{}, err
			}
			destination, err := strconv.Atoi(numbers[0])
			if err != nil {
				return [][]rangedMap{}, err
			}
			length, err := strconv.Atoi(numbers[2])
			if err != nil {
				return [][]rangedMap{}, err
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
	return allMaps, nil
}

// Learning Go: TODO s
// - understand packages, modules, imports, both when only locally and when shared packages online
// - understand defer keyword
