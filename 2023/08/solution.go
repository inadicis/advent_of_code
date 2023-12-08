package main

import (
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
	// GOLANG? Why is the type of x `any` and not `Way`?
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

func isStart(node string) bool {
	return node[len(node)-1] == 'A'
}

// func calculateNextExit(position Position, instructions string, nodesMap map[string][2]string) (NextExit, error) {
// 	// fmt.Printf("   Calculating next exit for %T, %v, ...\n", position, position)

// 	instructionI := position.instructionIndex
// 	currentNode := position.node
// 	var steps int
// 	for {
// 		steps++
// 		var instruction int
// 		if instructions[instructionI] == 'R' {
// 			instruction = 1
// 		} else if instructions[instructionI] != 'L' {
// 			return NextExit{}, fmt.Errorf("Instruction not parsable at i: %d", instructionI)
// 		}
// 		possibleNodes, exists := nodesMap[currentNode]
// 		if !exists {
// 			return NextExit{}, fmt.Errorf("No node found for %s", currentNode)
// 		}

// 		// fmt.Printf("Step %d, instructionI %d, instruction: %d, possibleNodes %v\n", steps, currentInstructionIndex, instruction, possibleNodes)
// 		nextNode := possibleNodes[instruction]
// 		// fmt.Printf("    i %d, currentNode %s\n", steps, currentNode)
// 		if nextNode == currentNode {
// 			return NextExit{}, fmt.Errorf("Entered Infinite loop while looking for exit for %T %v, final node: %s", position, position, nextNode)
// 		}
// 		currentNode = nextNode
// 		if isExit(currentNode) {
// 			// fmt.Printf("   nextNode %v is exit -> break\n", nextNode)
// 			break
// 		}
// 		instructionI++
// 		if instructionI >= len(instructions) {
// 			instructionI = instructionI % len(instructions)
// 		}
// 	}
// 	return NextExit{distance: steps, exit: currentNode, origin: position.node}, nil
// }

// type NextExit struct {
// 	origin   string
// 	distance int
// 	exit     string
// }

// type Position struct {
// 	instructionIndex int
// 	node             string
// }

// func InitCachedGetNextExit(instructions string, nodesMap map[string][2]string) func(position Position) (NextExit, error) {
// 	var cache = make(map[Position]NextExit)
// 	return func(position Position) (NextExit, error) {
// 		if retrieved, ok := cache[position]; ok {
// 			// fmt.Print(" (cached)")
// 			return retrieved, nil
// 		}
// 		r, err := calculateNextExit(position, instructions, nodesMap)
// 		if err != nil {
// 			return NextExit{}, err
// 		}
// 		cache[position] = r
// 		return r, nil
// 	}

// }

// func allSameDistance(waysHeap *WayHeap) bool {
// 	d := (*waysHeap)[0].distance
// 	// if d == 0 {
// 	// 	return false
// 	// }
// 	// fmt.Printf("Heap: [")
// 	for _, way := range *waysHeap {
// 		// fmt.Printf(" %v, ", way)
// 		if way.distance != d {
// 			// fmt.Printf("...]\n")
// 			return false
// 		}
// 	}
// 	// fmt.Printf("]\n")
// 	return true
// }

type Cycle struct {
	position         int
	length           int
	node             string
	nextInstructions string
	// exits            []Exit
	exitsInCycle     []Exit
	exitsBeforeCycle []Exit
}

type Position struct {
	nextInstructions string // instead of saving only the index, we save the concatenation to detect even tinier cycles
	node             string
	// todo careful about cycles where L/R are not relevant, could be tinier than here
	// if destination[XXX][0] == destination[XXX][1]
}

type Exit struct {
	node     string
	distance int
}

func findCycle(startNode string, instructions string, directions map[string][2]string) (Cycle, error) {
	var i int
	currentPosition := Position{
		node:             startNode,
		nextInstructions: instructions,
	}
	currentInstructionIndex := 0

	exits := []Exit{}
	visitedPositions := make(map[Position]int)
	for {
		if isExit(currentPosition.node) {
			exits = append(exits, Exit{node: currentPosition.node, distance: i})
		}
		startIndex, visited := visitedPositions[currentPosition]
		if visited {
			firstExitInCycleIndex := 0
			for i, exit := range exits {
				if exit.distance >= startIndex {
					firstExitInCycleIndex = i
					break
				}
			}
			return Cycle{
				position:         startIndex,
				length:           i - startIndex,
				node:             currentPosition.node,
				nextInstructions: currentPosition.nextInstructions,
				exitsInCycle:     exits[firstExitInCycleIndex:],
				exitsBeforeCycle: exits[:firstExitInCycleIndex],
			}, nil
		}
		visitedPositions[currentPosition] = i
		i++
		currentInstructionIndex++
		if currentInstructionIndex == len(instructions) {
			// could use >= or modulus but this should be fastest
			currentInstructionIndex = 0
		}
		// nextNode :=
		var instruction int

		if instructions[currentInstructionIndex] == 'R' {
			// equivalent to `if currentPosition.nextInstructions[0] == 'R'`
			instruction = 1
		} else if instructions[currentInstructionIndex] != 'L' {
			return Cycle{}, fmt.Errorf("Instruction not parsable at i: %d", currentInstructionIndex)
		}
		possibleNodes, exists := directions[currentPosition.node]
		if !exists {
			return Cycle{}, fmt.Errorf("No node found for %s", currentPosition.node)
		}
		nextNode := possibleNodes[instruction]

		currentPosition = Position{
			node:             nextNode,
			nextInstructions: instructions[currentInstructionIndex:] + instructions[:currentInstructionIndex],
		}
	}
}

func getResultPart2(filename string) (int, error) {
	text, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	instructions, nodesBlock, found := strings.Cut(strings.TrimSpace(string(text)), "\r\n\r\n")
	if !found {
		return 0, errors.New("Could not split file into instructions / nodes")
	}

	// fmt.Println(instructionsLine)
	// fmt.Println("----")
	// fmt.Println(nodesBlock)
	startNodes := []string{}
	directions := make(map[string][2]string)
	for i, nodeLine := range strings.Split(nodesBlock, "\r\n") {
		source, destinations, found := strings.Cut(nodeLine, " = ")
		if !found {
			return 0, fmt.Errorf("Could not split line i=%d, = not found.\n line: %s", i, nodeLine)
		}
		// fmt.Printf("i %d, source %s, destinations %s \n", i, source, destinations)
		if isStart(source) {
			startNodes = append(startNodes, source)
			// fmt.Printf(" found new start node. Current start nodes: %v\n", startNodes)
		}
		directions[source] = [2]string{destinations[1:4], destinations[6:9]}
	}

	// getNextExit := InitCachedGetNextExit(instructionsLine, nodes)

	// fmt.Printf("startNodes %v\n", startNodes)
	// currentWays := make([]Way, len(startNodes))
	// waysHeap := make(WayHeap, len(startNodes))

	foundCycles := make([]Cycle, len(startNodes))
	for i, node := range startNodes {
		// currentWays[i] = Way{distance: 0, position: Position{instructionIndex: 0, node: node}}
		// waysHeap[i] = &Way{distance: 0, position: Position{instructionIndex: 0, node: node}}
		foundCycles[i], err = findCycle(node, instructions, directions)
	}
	fmt.Println()
	fmt.Printf("found cycles %#v", foundCycles)
	return 0, nil

	// heap.Init(&waysHeap)

	// var minWay *Way
	// counter := 0
	// // fmt.Printf("Heap at start: %#v\n", waysHeap)
	// for {
	// 	counter++
	// 	// if counter > 10 {
	// 	// 	break
	// 	// }
	// 	if counter%10000000 == 0 {
	// 		fmt.Printf(" i: %d\n", counter)
	// 	}
	// 	minWay = heap.Pop(&waysHeap).(*Way) // Type assertion
	// 	// fmt.Printf("  Processing %T %#v , %v ...", minWay, minWay, minWay)
	// 	nextExit, err := getNextExit(minWay.position)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	// fmt.Printf("  found next %T %#v , %v \n", nextExit, nextExit, nextExit)
	// 	newDistance := minWay.distance + nextExit.distance
	// 	heap.Push(&waysHeap, &Way{
	// 		distance: newDistance,
	// 		position: Position{
	// 			instructionIndex: newDistance % len(instructions),
	// 			node:             nextExit.exit,
	// 		}},
	// 	)

	// 	if allSameDistance(&waysHeap) {
	// 		break
	// 	}
	// }
	// finalWay := heap.Pop(&waysHeap).(*Way)
	// fmt.Printf("Calculation Completed after %d iterations. distance: %d\n", counter, finalWay.distance)
	// return finalWay.distance, nil
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
