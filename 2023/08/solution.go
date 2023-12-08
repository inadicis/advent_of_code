package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

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
	var instruction int
	if instructions[position.instructionIndex] == 'R' {
		instruction = 1
	} else if instructions[position.instructionIndex] != 'L' {
		return NextExit{}, fmt.Errorf("Instruction not parsable at i: %d", position.instructionIndex)
	}

	currentNode := position.node
	var steps int
	for ; !isExit(currentNode); steps++ {
		// fmt.Printf("  i %d, currentNode %s\n", i, currentNode)
		possibleNodes, exists := nodesMap[currentNode]
		if !exists {
			return NextExit{}, fmt.Errorf("No node found for %s", currentNode)
		}

		// fmt.Printf("Step %d, instructionI %d, instruction: %d, possibleNodes %v\n", steps, currentInstructionIndex, instruction, possibleNodes)
		currentNode = possibleNodes[instruction]
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
			fmt.Println("from cache")
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

	// fmt.Printf("startNodes %v\n", startNodes)
	currentPositions := make([]Position, len(startNodes))
	for i, node := range startNodes {
		currentPositions[i] = Position{instructionIndex: 0, node: node}
	}
	amountCopied := copy(currentNodes, startNodes) // why does copy not work?
	fmt.Printf(" amount copied: %d\n", amountCopied)
	// fmt.Printf("currentNodes %v\n", currentNodes)
	var currentInstructionIndex int
	var instruction int
	var steps int
	for {
		looksGood := true
		if instructionsLine[currentInstructionIndex] == 'R' {
			instruction = 1
		} else if instructionsLine[currentInstructionIndex] == 'L' {
			instruction = 0 // could be removed but I want the defensive last else
		} else {
			return 0, fmt.Errorf("Instruction not parsable at i: %d", currentInstructionIndex)
		}
		steps++
		if steps%1000000 == 0 {
			fmt.Printf("step %d, currentNodes %v\n", steps, currentNodes)
		}
		for i, currentNode := range currentNodes {
			// fmt.Printf("  i %d, currentNode %s\n", i, currentNode)
			possibleNodes, exists := nodes[currentNode]
			if !exists {
				return 0, fmt.Errorf("No node found for %s", currentNode)
			}

			// fmt.Printf("Step %d, instructionI %d, instruction: %d, possibleNodes %v\n", steps, currentInstructionIndex, instruction, possibleNodes)
			currentNode = possibleNodes[instruction]

			if currentNode[len(currentNode)-1] != 'Z' {
				looksGood = false
			}
			currentNodes[i] = currentNode
		}
		if looksGood {
			return steps, nil
		}

		currentInstructionIndex++
		if currentInstructionIndex >= len(instructionsLine) {
			currentInstructionIndex = currentInstructionIndex % len(instructionsLine)
		}

	}

}
