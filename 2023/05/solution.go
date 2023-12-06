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

func maxInt() int {
	return 10000000000
}

func main() {
	result, err := solution5("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("result: %d", result)
	// mappings := []rangedMap{
	// 	{source_start: 0, destination_start: 10, length: 5},
	// 	// rangedMap{source_start: 5, destination_start: 5, length: 5},
	// 	{source_start: 10, destination_start: 0, length: 5},
	// 	{source_start: 15, destination_start: 15, length: 3},
	// 	// {source_start: 18, destination_start: 30, length: 4},
	// 	// {source_start: 25, destination_start: 18, length: 4},
	// }
	// firstMappings := fillGaps(mappings)
	// secondMappings := fillGaps([]rangedMap{
	// 	{source_start: 0, destination_start: 7, length: 8},
	// 	// rangedMap{source_start: 5, destination_start: 5, length: 5},
	// 	{source_start: 8, destination_start: 0, length: 15},
	// })
	// // n := 4
	// // match, index := findRelevantMap(n, mappings)
	// // fmt.Printf("Searched for n=%d, found map %+v at position %d", n, match, index)
	// // filledMappings := fillGaps(mappings)
	// // fmt.Printf("result: %+v", filledMappings)

	// compressedMappings := compressRangedMaps(firstMappings, secondMappings)
	// fmt.Printf("result: %+v", compressedMappings)

}

func compressRangedMaps(firstMaps []rangedMap, secondMaps []rangedMap) []rangedMap {
	// assumes both slices are sorted (by source_start), with no gap (see fillGaps) and no overlap
	compressedMap := []rangedMap{}

	for _, firstMap := range firstMaps {
		fmt.Printf("first Map: %+v\n", firstMap)
		_, minIndex := findRelevantMap(firstMap.destination_start, secondMaps)
		// fmt.Printf("  Relevant second Minimum maps: %+v\n", minMap)
		_, maxIndex := findRelevantMap(firstMap.destination_start+firstMap.length-1, secondMaps)
		// fmt.Printf("  Relevant second Maximum maps: %+v\n", maxMap)
		relevantSecondMaps := secondMaps[minIndex : maxIndex+1]
		fmt.Printf("  Relevant second maps (%d): %+v\n", len(relevantSecondMaps), relevantSecondMaps)

		current_start := firstMap.source_start
		remaining := firstMap.length
		for _, secondMap := range relevantSecondMaps { // TODO careful if second ranges are more large than first map (min, max around!)
			startOffset := firstMap.destination_start - secondMap.source_start
			length := secondMap.length
			destination_start := secondMap.destination_start
			if startOffset > 0 { // can only happen with the first of relevantSecondMaps
				destination_start += startOffset
				length -= startOffset
			}

			length = min(remaining, length) // avoid endOffset
			newRange := rangedMap{source_start: current_start, destination_start: destination_start, length: length}
			compressedMap = append(compressedMap, newRange)
			current_start += length
			remaining -= length
			fmt.Printf("  new range: %+v\n", newRange)
		}

	}
	return compressedMap
}

func fillGaps(mappings []rangedMap) []rangedMap {
	// make all mappings explicit, by having even mappings 1:1
	fullMapping := []rangedMap{}

	current_start := 0
	for _, m := range mappings {
		if m.source_start != current_start {
			newRangeMap := rangedMap{source_start: current_start, destination_start: current_start, length: m.source_start - current_start}
			fullMapping = append(fullMapping, newRangeMap)
		}
		fullMapping = append(fullMapping, m)
		current_start = m.source_start + m.length
	}
	fullMapping = append(fullMapping, rangedMap{source_start: current_start, destination_start: current_start, length: maxInt() - current_start})
	// TODO length, max accepted size of n, represents positive infinite
	return fullMapping
}

func findRelevantMap(n int, rangedMaps []rangedMap) (rangedMap, int) {
	// assumes slice is sorted (by source_start), with no gap (see fillGaps) and no overlap
	// binary search
	l := 0
	r := len(rangedMaps) - 1
	for {
		i := l + (r-l)/2
		currentMap := rangedMaps[i]
		if n < currentMap.source_start {
			r = i - 1
			continue
		}
		if n >= currentMap.source_start+currentMap.length {
			l = i + 1
			continue
		}
		return currentMap, i
	}

}

/*
Brute force solution is now implemented but will take circa 15 minutes to run
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
	filledMaps := [][]rangedMap{}
	for _, mappings := range allMaps {
		filledMaps = append(filledMaps, fillGaps(mappings))
	}
	finalMappings := filledMaps[0]
	for _, rangedMap := range filledMaps[1:] {
		finalMappings = compressRangedMaps(finalMappings, rangedMap)
	}

	locations := []int{}
	for _, seedRange := range seedRanges {
		fmt.Printf("seed range %+v\n", seedRange)
		for i := 0; i < seedRange.length; i++ {
			seed := seedRange.start + i
			rm, _ := findRelevantMap(seed, finalMappings)
			locations = append(locations, rm.destination_start+seed-rm.source_start)
		}

	}

	minLocations := []int{}

	for _, seedRange := range seedRanges {
		// TODO duplicated with compress function
		_, minIndex := findRelevantMap(seedRange.start, finalMappings)
		_, maxIndex := findRelevantMap(seedRange.start+seedRange.length, finalMappings)
		relevantMaps := finalMappings[minIndex : maxIndex+1]
		// we only need the minimum, so we can simply check the first value for each of these relevant maps

		current_start := seedRange.start
		remaining := seedRange.length
		for _, secondMap := range relevantMaps {
			startOffset := seedRange.start - secondMap.source_start
			length := secondMap.length
			destination_start := secondMap.destination_start
			if startOffset > 0 {
				destination_start += startOffset
				length -= startOffset
			}
			minLocations = append(minLocations, destination_start)

			length = min(remaining, length) // avoid endOffset
			current_start += length
			remaining -= length
		}

	}
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
