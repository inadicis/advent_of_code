package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
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
	start := time.Now()
	result, err := solution5("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	fmt.Printf("result: %d, calculated in %d ms", result, end.Sub(start).Microseconds())

}

func compressRangedMaps(firstMaps []rangedMap, secondMaps []rangedMap) []rangedMap {
	// assumes both slices are sorted (by source_start), with no gap (see fillGaps) and no overlap
	compressedMap := []rangedMap{}

	for _, firstMap := range firstMaps {
		_, minIndex := findRelevantMap(firstMap.destination_start, secondMaps)
		_, maxIndex := findRelevantMap(firstMap.destination_start+firstMap.length-1, secondMaps)
		relevantSecondMaps := secondMaps[minIndex : maxIndex+1]

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
	return slices.Min(minLocations), nil
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
