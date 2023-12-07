package main

import (
	"fmt"
	"log"
)

func main() {
	r, err := getResult("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}

func getResult(filename string) (int, error) {

	return 0
}
