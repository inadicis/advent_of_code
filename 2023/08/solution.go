package main

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// An IntHeap is a min-heap of ints.
type Way struct {
	position Position
	distance int
}

type WayHeap []*Way

func (h WayHeap) Len() int           { return len(h) }
func (h WayHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h WayHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *WayHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Way))
}

func (h *WayHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	r, err := getResultPart2("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}

func getResultPart1(filename string) (int, error) {
	text, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	instructionsLine, nodesBlock, found := strings.Cut(strings.TrimSpace(string(text)), "\r\n\r\n")
	if !found {
		return 0, errors.New("Could not split file into instructions / nodes")
	}

	// fmt.Println(instructionsLine)
	// fmt.Println("----")
	// fmt.Println(nodesBlock)

	nodes := make(map[string][2]string)
	for i, nodeLine := range strings.Split(nodesBlock, "\r\n") {
		source, destinations, found := strings.Cut(nodeLine, " = ")
		if !found {
			return 0, fmt.Errorf("Could not split line i=%d, = not found.\n line: %s", i, nodeLine)
		}
		nodes[source] = [2]string{destinations[1:4], destinations[6:9]}
	}

	currentNode := "AAA"
	var currentInstructionIndex int
	var instruction int
	var steps int
	for {
		if instructionsLine[currentInstructionIndex] == 'R' {
			instruction = 1
		} else if instructionsLine[currentInstructionIndex] == 'L' {
			instruction = 0 // could be removed but I want the defensive last else
		} else {
			return 0, fmt.Errorf("Instruction not parsable at i: %d", currentInstructionIndex)
		}
		possibleNodes, exists := nodes[currentNode]
		if !exists {
			return 0, fmt.Errorf("No node found for %s", currentNode)
		}

		// fmt.Printf("Step %d, instructionI %d, instruction: %d, possibleNodes %v", steps, currentInstructionIndex, instruction, possibleNodes)
		currentNode = possibleNodes[instruction]
		steps++

		if currentNode == "ZZZ" {
			return steps, nil
		}

		currentInstructionIndex++
		if currentInstructionIndex >= len(instructionsLine) {
			currentInstructionIndex = currentInstructionIndex % len(instructionsLine)
		}

	}

}

func isExit(node string) bool {
	return node[len(node)-1] == 'Z'
}

func calculateNextExit(position Position, instructions string, nodesMap map[string][2]string) (NextExit, error) {
	fmt.Printf("   Calculating next exit for %T, %v, ...\n", position, position)

	instructionI := position.instructionIndex
	currentNode := position.node
	var steps int
	for {
		steps++
		var instruction int
		if instructions[instructionI] == 'R' {
			instruction = 1
		} else if instructions[instructionI] != 'L' {
			return NextExit{}, fmt.Errorf("Instruction not parsable at i: %d", instructionI)
		}
		possibleNodes, exists := nodesMap[currentNode]
		if !exists {
			return NextExit{}, fmt.Errorf("No node found for %s", currentNode)
		}

		// fmt.Printf("Step %d, instructionI %d, instruction: %d, possibleNodes %v\n", steps, currentInstructionIndex, instruction, possibleNodes)
		nextNode := possibleNodes[instruction]
		fmt.Printf("    i %d, currentNode %s\n", steps, currentNode)
		if nextNode == currentNode {
			return NextExit{}, fmt.Errorf("Entered Infinite loop while looking for exit for %T %v, final node: %s", position, position, nextNode)
		}
		currentNode = nextNode
		if isExit(currentNode) {
			fmt.Printf("   nextNode %v is exit -> break\n", nextNode)
			break
		}
		instructionI++
		if instructionI >= len(instructions) {
			instructionI = instructionI % len(instructions)
		}
	}
	return NextExit{distance: steps, node: currentNode}, nil
}

type NextExit struct {
	distance int
	node     string
}

type Position struct {
	instructionIndex int
	node             string
}

func InitCachedGetNextExit(instructions string, nodesMap map[string][2]string) func(position Position) (NextExit, error) {
	var cache = make(map[Position]NextExit)
	return func(position Position) (NextExit, error) {
		if retrieved, ok := cache[position]; ok {
			fmt.Print(" (cached)")
			return retrieved, nil
		}
		r, err := calculateNextExit(position, instructions, nodesMap)
		if err != nil {
			return NextExit{}, err
		}
		cache[position] = r
		return r, nil
	}

}

func allSameDistance(waysHeap *WayHeap) bool {
	d := (*waysHeap)[0].distance
	// if d == 0 {
	// 	return false
	// }
	fmt.Printf("Heap: [")
	for _, way := range *waysHeap {
		fmt.Printf(" %v, ", way)
		if way.distance != d {
			fmt.Printf("...]\n")
			return false
		}
	}
	fmt.Printf("]\n")
	return true
}

func getResultPart2(filename string) (int, error) {
	// needs at least minutes to find solutions
	// optimization possibilities:
	// - find cycles for each source, for each cycle the index of possible destinations (end with Z)
	// 	 once all cycles found -> calculate equation solution
	// - define cached function )nextExit(sourceNode string, instructions string, index int)(steps int, exit string
	text, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	instructionsLine, nodesBlock, found := strings.Cut(strings.TrimSpace(string(text)), "\r\n\r\n")
	if !found {
		return 0, errors.New("Could not split file into instructions / nodes")
	}

	// fmt.Println(instructionsLine)
	// fmt.Println("----")
	// fmt.Println(nodesBlock)
	startNodes := []string{}
	nodes := make(map[string][2]string)
	for i, nodeLine := range strings.Split(nodesBlock, "\r\n") {
		source, destinations, found := strings.Cut(nodeLine, " = ")
		if !found {
			return 0, fmt.Errorf("Could not split line i=%d, = not found.\n line: %s", i, nodeLine)
		}
		// fmt.Printf("i %d, source %s, destinations %s \n", i, source, destinations)
		if source[len(source)-1] == 'A' {
			startNodes = append(startNodes, source)
			// fmt.Printf(" found new start node. Current start nodes: %v\n", startNodes)
		}
		nodes[source] = [2]string{destinations[1:4], destinations[6:9]}
	}

	getNextExit := InitCachedGetNextExit(instructionsLine, nodes)

	fmt.Printf("startNodes %v\n", startNodes)
	// currentWays := make([]Way, len(startNodes))
	waysHeap := make(WayHeap, len(startNodes))

	for i, node := range startNodes {
		// currentWays[i] = Way{distance: 0, position: Position{instructionIndex: 0, node: node}}
		waysHeap[i] = &Way{distance: 0, position: Position{instructionIndex: 0, node: node}}
	}

	// h := WayHeap(currentWays)
	heap.Init(&waysHeap)
	// heap.Push(h, 3)
	// fmt.Printf("minimum: %d\n", (*h)[0])

	// for h.Len() > 0 {
	// 	fmt.Printf("%d ", heap.Pop(h))
	// }

	var minWay *Way
	counter := 0
	fmt.Printf("Heap at start: %#v\n", waysHeap)
	for {
		counter++
		// if counter > 10 {
		// 	break
		// }
		fmt.Printf(" i: %d", counter)
		minWay = heap.Pop(&waysHeap).(*Way) // Type assertion
		fmt.Printf("  Processing %T %#v , %v ...", minWay, minWay, minWay)
		nextExit, err := getNextExit(minWay.position)
		if err != nil {
			return 0, err
		}
		fmt.Printf("  found next %T %#v , %v \n", nextExit, nextExit, nextExit)
		newDistance := minWay.distance + nextExit.distance
		heap.Push(&waysHeap, &Way{
			distance: newDistance,
			position: Position{
				instructionIndex: newDistance % len(instructionsLine),
				node:             nextExit.node,
			}},
		)

		if allSameDistance(&waysHeap) {
			break
		}
	}
	finalWay := heap.Pop(&waysHeap).(*Way)
	fmt.Printf("Calculation Completed after %d iterations. distance: %d\n", counter, finalWay.distance)
	return finalWay.distance, nil
	// var currentInstructionIndex int
	// var instruction int
	// var steps int
	// for {
	// 	looksGood := true
	// 	if instructionsLine[currentInstructionIndex] == 'R' {
	// 		instruction = 1
	// 	} else if instructionsLine[currentInstructionIndex] == 'L' {
	// 		instruction = 0 // could be removed but I want the defensive last else
	// 	} else {
	// 		return 0, fmt.Errorf("Instruction not parsable at i: %d", currentInstructionIndex)
	// 	}
	// 	steps++
	// 	if steps%1000000 == 0 {
	// 		fmt.Printf("step %d, currentNodes %v\n", steps, currentNodes)
	// 	}
	// 	for i, currentNode := range currentNodes {
	// 		// fmt.Printf("  i %d, currentNode %s\n", i, currentNode)
	// 		possibleNodes, exists := nodes[currentNode]
	// 		if !exists {
	// 			return 0, fmt.Errorf("No node found for %s", currentNode)
	// 		}

	// 		// fmt.Printf("Step %d, instructionI %d, instruction: %d, possibleNodes %v\n", steps, currentInstructionIndex, instruction, possibleNodes)
	// 		currentNode = possibleNodes[instruction]

	// 		if currentNode[len(currentNode)-1] != 'Z' {
	// 			looksGood = false
	// 		}
	// 		currentNodes[i] = currentNode
	// 	}
	// 	if looksGood {
	// 		return steps, nil
	// 	}

	// 	currentInstructionIndex++
	// 	if currentInstructionIndex >= len(instructionsLine) {
	// 		currentInstructionIndex = currentInstructionIndex % len(instructionsLine)
	// 	}

	// }

}
