package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func pprint(boxes []list.List) {

	for i, box := range boxes {
		s := ""
		for e := box.Front(); e != nil; e = e.Next() {
			s += fmt.Sprintf("[%v]", e.Value)
		}
		fmt.Printf(" box %03d: %s\n", i, s)
	}
}

func main() {
	text, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	instructions := cleanupData(string(text))
	result := part2(instructions)
	fmt.Printf("\nFinal result: %d", result)

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

type Value struct {
	focal string
	label string
}

func part2(instructions []string) int {

	boxes := make([]list.List, 256)
	for _, instruction := range instructions {
		label, focalLength, insert := strings.Cut(instruction, "=")
		if !insert {
			label = instruction[:len(instruction)-1]
		}
		h := hash(label)
		inserted := false
		box := &boxes[h]
		for e := box.Front(); e != nil; e = e.Next() {
			if e.Value.(Value).label == label {
				if insert {
					e.Value = Value{label: label, focal: focalLength}
					inserted = true
				} else {
					box.Remove(e)
				}
				break
			}
		}
		if insert && !inserted {
			box.PushBack(Value{label: label, focal: focalLength})
		}

	}
	total := 0
	for i, box := range boxes {
		slot := 1
		for e := box.Front(); e != nil; e = e.Next() {
			f, err := strconv.Atoi(e.Value.(Value).focal)
			if err != nil {
				panic(err)
			}
			total += (i + 1) * slot * f
			slot++
		}
	}
	return total
}
