package utils

import (
	"bytes"
	"fmt"
	"strings"
)

func pprint(lines [][]rune) {
	for _, row := range lines {
		for _, char := range row {
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
}

func same(a [][]rune, b [][]rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, row := range a {
		if len(row) != len(b[i]) {
			return false
		}
		for j, c := range row {
			if c != b[i][j] {
				return false
			}
		}
	}
	return true
}

func encode(a [][]rune) string {
	var buffer bytes.Buffer
	for _, row := range a {
		for _, char := range row {
			buffer.WriteRune(char)
		}
		buffer.WriteRune('\n')
	}
	return buffer.String()
}

func decode(s string) [][]rune {
	splitLines := strings.Split(s, "\n")
	a := make([][]rune, len(splitLines))
	for i, l := range splitLines {
		r := make([]rune, len(l))
		for j, char := range l {
			r[j] = char
		}
		a[i] = r
	}
	return a
}

func StringToRunes(s string) []rune {
	r := make([]rune, len(s))
	for i, c := range s {
		r[i] = c
	}
	return r
}

func StringsToRunes(a []string) [][]rune {
	rows := make([][]rune, len(a))
	for i, s := range a {
		rows[i] = StringToRunes(s)
	}
	return rows
}
