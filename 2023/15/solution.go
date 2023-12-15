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

	// boxes := make([]list.List, 10)
	// boxes[0].PushBack(5)
	// pprint(boxes)
	// boxes[3].PushBack(4)
	// boxes[3].PushBack(10)
	// pprint(boxes)
	// box := boxes[1]
	// box.PushBack(1000)
	// pprint(boxes)
	text, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	instructions := cleanupData(string(text))
	result := part2(instructions)
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

type Value struct {
	focal string
	label string
}

func part2(instructions []string) int {

	boxes := make([]list.List, 256)
	for _, instruction := range instructions {
		fmt.Printf(" working with instruction %q\n", instruction)
		label, focalLength, insert := strings.Cut(instruction, "=")
		if !insert {
			label = instruction[:len(instruction)-1]
		}
		h := hash(label)
		inserted := false
		fmt.Printf(" label %q, insert %v, focal %s, hash %d\n", label, insert, focalLength, h)
		for e := boxes[h].Front(); e != nil; e = e.Next() {
			fmt.Printf("  slot -> %v\n", e)
			if e.Value.(Value).label == label {
				if insert {
					e.Value = Value{label: label, focal: focalLength}
					inserted = true
				} else {
					boxes[h].Remove(e)
				}
				fmt.Printf("  box after change: %v\n", boxes[h])
				break
			}
		}
		if insert && !inserted {
			fmt.Printf("  inserting %q at label %q, hash %d\n", focalLength, label, h)
			boxes[h].PushBack(Value{label: label, focal: focalLength})
		}
		for e := boxes[h].Front(); e != nil; e = e.Next() {
			fmt.Printf(" | label %q, focal %q", e.Value.(Value).label, e.Value.(Value).focal)
		}
		fmt.Println()

	}
	for i, box := range boxes {
		s := ""
		for e := box.Front(); e != nil; e = e.Next() {
			s += fmt.Sprintf("[%v]", e.Value)
		}
		fmt.Printf(" box %03d: %s\n", i, s)
	}
	// fmt.Printf("%v", boxes)
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
