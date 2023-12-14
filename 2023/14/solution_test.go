package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParts(t *testing.T) {
	tests := []struct {
		expected int
		filename string
		fn       func(string) int
	}{
		{
			expected: 136,
			filename: "test.txt",
			fn:       part1,
		},
	}
	for _, test := range tests {
		b, err := os.ReadFile(test.filename)
		assert.NoError(t, err, test.filename)
		assert.Equal(t, test.expected, test.fn(string(b)))
	}
}

// func TestInternalFunc(t *testing.T) {
// 	tests := []struct {
// 		expected int64
// 		input    string
// 		fn       func(string) int64
// 	}{
// 		{
// 			expected: 999,
// 			input:    "1,2,4,5 A B C",
// 			fn:       part1,
// 		},
// 	}
// 	for _, test := range tests {
// 		assert.Equal(t, test.expected, test.fn(test.input))
// 	}
// }
