package main

import (
	"regexp"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello("Gladys")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}

var cardsTestData = []struct {
	in  string
	out int
}{
	{"Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53", 8},
	{"Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19", 2},
	{"Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1", 2},
	{"Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83", 1},
	{"Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36", 0},
	{"Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11", 0},
}

func TestGetCardPoints(t *testing.T) {
	for _, tt := range cardsTestData {
		amount, err := GetCardPoints(tt.in)
		if err != nil {
			t.Errorf("GetCardPoints(%q) errored: %s", tt.in, err)
		}
		if amount != tt.out {
			t.Errorf("GetCardPoints(%q) expected %d, actual %d", tt.in, tt.out, amount)
		}
	}
}
