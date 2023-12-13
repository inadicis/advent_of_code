package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nFinal result: %d", part1(string(b)))
}

func part1(input string) int64 {
	data := cleanupInput(input)
	for _, d := range data {
		fmt.Println(d)
	}
	return -1
}

func part2(input string) int64 {
	data := cleanupInput(input)
	for _, d := range data {
		fmt.Println(d)
	}
	return -1
}

func cleanupInput(input string) []string {
	return strings.Split(input, "\r\n")
}
