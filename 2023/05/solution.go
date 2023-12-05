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
	// fmt.Println(seeds)

	// sort.Sort()
	// slices.SortFunc()
	allMaps := [][]rangedMap{}
	for i, block := range text_blocks[1:] {
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
		for _, m := range rangedMaps {
			fmt.Printf("%d | source_start: %d, destination_start: %d, length: %d \n", i, m.source_start, m.destination_start, m.length)
		}
		fmt.Println("---------------")
		allMaps = append(allMaps, rangedMaps)
	}
	locations := []int{}
	for i, seed := range seeds {
		fmt.Println(i, seed)
		current, err := strconv.Atoi(seed)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		for j, rangedMaps := range allMaps {
			// linear search of the correct range. A binary search should be possible though
			fmt.Printf(" Step %d, Maps: %+v\n", j, rangedMaps)
			for _, rangedMap := range rangedMaps {
				// if rangedMap.source_start > current {
				// 	continue // not
				// }

				diff := current - rangedMap.source_start
				fmt.Printf("   current: %d, start: %d, length: %d, diff: %d\n", current, rangedMap.source_start, rangedMap.length, diff)

				if diff < 0 {
					fmt.Printf("  No range found, current n remains: %d\n", current)
					break // if there is a map, it is further down (because sorted)
				}
				if diff < rangedMap.length {
					current = rangedMap.destination_start + diff
					fmt.Printf("  found range %+v, new current n: %d\n", rangedMap, current)
					break // map found, go to next step
				} //else {
				fmt.Printf("    Not yet, continue\n")
				//continue // no map existing for this number (assuming no overlaps)
				//}
			}
			fmt.Printf(" n=%d", current)
		}
		fmt.Printf("\n Result: %d", current)
		locations = append(locations, current)
	}
	// fmt.Println(locations)
	min_location := locations[0]
	for _, location := range locations[1:] {
		if location < min_location {
			min_location = location
		}
	}
	return min_location
	// return min(locations...) -> TODO why does that not work? min wants fixed number of parameters and not variatic????

	// for _, s := range text {
	// 	fmt.Println(s)
	// 	fmt.Println("--------------------------")
	// }

	// fmt.Println(text)

	// file, err := os.Open("input.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)
	// var seeds []string
	// maps := make([][]rangedMap, 7)
	// skip_next := false
	// current_map := 0

	// cleanup maps
	// for i := 0; scanner.Scan(); i++ { // each line
	// 	current_map := []rangedMap{}
	// 	line := scanner.Text()

	// 	if line == "" {
	// 		current_map += 1
	// 		continue
	// 	}
	// 	fmt.Print(line)
	// }

	// if err := scanner.Err(); err != nil {
	// 	fmt.Println(err)
	// }

}

// Learning Go: TODO s
// - understand packages, modules, imports, both when only locally and when shared packages online
// - understand defer keyword
