package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	text, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	instructions := cleanupData(string(text))
	result := part1(instructions)
	fmt.Printf("\nFinal result: %d", result)
	// fmt.Print(hash("HASH"))
}
func cleanupData(text string) []string {
	return strings.Split(text, ",")
}

func hash(s string) int {
	h := 0
	for _, r := range s {
		h = ((h + int(r)) * 17) % 256
	}
	return h
}

func part1(instructions []string) int {
	total := 0
	for _, instruction := range instructions {
		total += hash(instruction)
	}
	return total
}
