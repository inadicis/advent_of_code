package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
)

// go mod tidy to check dependencies

func main() {
	// fmt.Println("Hello, World!")

	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	//log.SetPrefix("greetings: ")
	log.SetFlags(0)

	// Request a greeting message.
	message, err := Hello("Gladys")
	// If an error was returned, print it to the console and
	// exit the program.
	if err != nil {
		log.Fatal(err)
	}

	// If no error was returned, print the returned message
	// to the console.
	fmt.Println(message)

	// ioutil.ReadFile(filename)

}

func Hello(name string) (string, error) {

	if name == "" {
		return "", errors.New("empty name")
	}

	message := fmt.Sprintf(randomFormat(), name)
	// := declares and initializes a variable, inducing its type . equivalent to:
	// var message string
	// message2 = fmt.Sprintf("Hello, %v, welcome!", name)
	// unused variables  -> error in go

	return message, nil
}

// Hellos returns a map that associates each of the named people
// with a greeting message.
func Hellos(names []string) (map[string]string, error) {
	// A map to associate names with messages.
	messages := make(map[string]string)
	// Loop through the received slice of names, calling
	// the Hello function to get a message for each name.
	for _, name := range names {
		message, err := Hello(name)
		if err != nil {
			return nil, err
		}
		// In the map, associate the retrieved message with
		// the name.
		messages[name] = message
	}
	return messages, nil
}

// randomFormat returns one of a set of greeting messages. The returned
// message is selected at random.
func randomFormat() string {
	// A slice of message formats.
	formats := []string{
		"Hi, %v. Welcome!",
		"Great to see you, %v!",
		"Hail, %v! Well met!",
	}
	return formats[rand.Intn(len(formats))]
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
*/
