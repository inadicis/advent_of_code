package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// go mod tidy to check dependencies

func main() {
	// fmt.Println("Hello, World!")

	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	//log.SetPrefix("greetings: ")
	log.SetFlags(0)
	content, err := os.ReadFile("input.txt")
	// use a scanner next time to not hold everything in memory
	if err != nil {
		log.Fatal(err)
	}
	text := strings.Trim(string(content), " \r\n") //strip newlines and spaces
	lines := strings.Split(text, "\r\n")
	counters := make([]int, len(lines))
	fmt.Println(len(counters))

	for _, line := range lines {
		CountCardCopies(line, counters)
		if err != nil {
			log.Fatal(err)
		}
	}
	totalCards := 0
	for _, amount := range counters {
		totalCards += amount
	}
	fmt.Println("result: ", totalCards)

}

func CountCardCopies(line string, counters []int) {

}

func GetCardPoints(line string) (int, error) {
	amount, _ := amountWinningMatches(line)
	if amount == 0 {
		return 0, nil
	}

	points := 1
	for i := 0; i < amount-1; i++ {
		points *= 2
	}
	return points, nil
}

func amountWinningMatches(line string) (int, error) {
	fmt.Println("getting points from line", line)
	_, numbers, found := strings.Cut(line, ":")
	if !found {
		log.Fatal("Not found :")
	}
	winning_numbers, actual_numbers, found := strings.Cut(numbers, "|")

	//numbers := strings.Split(line[len("Card 213:"):], "|")
	//winning_numbers, actual_numbers := strings.Trim(numbers[0]), strings.Trim(numbers[1])
	//strings.strip

	if !found {
		log.Fatal("Not found |")
	}
	winning := map[int]bool{}
	for _, number := range strings.Split(strings.Trim(winning_numbers, " "), " ") {
		if number == "" {
			continue
		}
		n, err := strconv.Atoi(number)
		if err != nil {
			fmt.Println("Could not recognize int from string", number)
			return 0, err
		}
		winning[n] = true
	}
	matches := []int{}

	for _, number := range strings.Split(strings.Trim(actual_numbers, " "), " ") {
		if number == "" {
			continue
		}
		n, err := strconv.Atoi(number)
		if err != nil {
			fmt.Println("Could not recognize int from string", number)
			return 0, err
		}

		if winning[n] {
			matches = append(matches, n)
		}
	}
	return len(matches), nil
}

/*
Notes about Go

	// fuction starting with capital letter -> usable outside this package: "exported name"

	// errors are not raised, but returned
	functions can return multiple values. They are not really returning a "tuple" of values -> when assigning ``m := f()``  -> m only takes the first returned value and drops the rest. `` _, n := f()``  would take the second one.
	i, j := 2 -> both equal 2 ?
	// break and continue allow labeling which loop / switch you want to break out from

	// slice is the equivalent of a list in python?
	//  -> No. it is an abstraction on top an actual (fixed length) Array.
	//     It allows slicing the array but not growing it outside its capacity. -> can grow slice buy building new array, copy it, write our own append function, or use builtin append()
	//     unpack operator "..." for slices, eg: ``a = append(a, 1, 2, b...)
	// curly braces to initialize + fill a slice?
	// what is the comparison to keyword "new" and keyword "make"?
	// map is equivalent of python dict?
	//   - map is a reference type, that's why declaring it ``var m map[string]int`` is fine for reads (empty) but not for writes -> no where to write to -> needs initiliazation ``m = make(map[string]int)``
	//   or ``m = map[string]int{}`` (empty map literal is equivalent to using ``make`` function)
	//   - accessing a map m["mykey"] returns the zero value of the value type (0 for int), but you can differentiate it from an actual 0 with the second returned value (boolean) `` value, ok := m["mykey"] ``

		// range is what? not the syntax of a function. and it returns index and value? or only value?
		// range keyword returns key, value for a map

	String Formatting with %X -> have to learn all these X
	- couldn't find what %q is, but used it for a struct

	"defer" keyword ?
*/
