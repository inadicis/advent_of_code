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

func getResultPart2(filename string) (int, error) {
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
		fmt.Printf("i %d, source %s, destinations %s \n", i, source, destinations)
		if source[len(source)-1] == 'A' {
			startNodes = append(startNodes, source)
			fmt.Printf(" found new start node. Current start nodes: %v\n", startNodes)
		}
		nodes[source] = [2]string{destinations[1:4], destinations[6:9]}
	}

	var currentNodes []string
	copy(currentNodes, startNodes) // why does copy not work?
	fmt.Printf("currentNodes %v\n", currentNodes)
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

		fmt.Printf("step %d, currentNodes %v", steps, currentNodes)
		for i, currentNode := range currentNodes {
			fmt.Printf("  i %d, currentNode %s", i, currentNode)
			possibleNodes, exists := nodes[currentNode]
			if !exists {
				return 0, fmt.Errorf("No node found for %s", currentNode)
			}

			fmt.Printf("Step %d, instructionI %d, instruction: %d, possibleNodes %v", steps, currentInstructionIndex, instruction, possibleNodes)
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
