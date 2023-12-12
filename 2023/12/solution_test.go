package main

import (
	"testing"
)

type TestCase struct {
	mask               string
	damagedGroupsSizes []int
	expected           int
}

func TestCountPossibilities(t *testing.T) {
	testData := []TestCase{
		{"", []int{1}, 0},
		{"#", []int{1}, 1},
		{"##", []int{1}, 0},
		{"?", []int{1}, 1},
		{"??", []int{1}, 2},
		{"???", []int{1}, 3},
		{"???...?", []int{1}, 4},
		{"???.#.?", []int{1}, 1},
		{".", []int{1}, 0},
		{"#?", []int{1, 1}, 0},
		{".#?.", []int{1, 1}, 0},
		{".#.?.", []int{1, 1}, 1},
		{".##.?.", []int{1, 1}, 0},
		{"???", []int{1}, 3},
		{"???", []int{2}, 2},
		{"?", []int{1}, 1},
		{"???.", []int{1}, 3},
		{".???.", []int{1}, 3},
		{"#???.", []int{1}, 1},
		{"#???.", []int{1, 2}, 1},
		{"?#??#", []int{2, 1}, 2},
		{"?#?#?", []int{1}, 0},
		{"?#??#???....#..", []int{1, 2}, 0},
	}
	for i, tCase := range testData {
		actual := countPossibilities(tCase.mask, tCase.damagedGroupsSizes)
		if actual != tCase.expected {
			t.Errorf("%02d (%q, %v): expected %d, actual %d", i, tCase.mask, tCase.damagedGroupsSizes, tCase.expected, actual)
		}
		// else {
		// 	fmt.Printf("%02d: expected = actual = %d\n", i, actual)
		// }
	}
}

???.###????.###????.###????.###????.### 1,1,3,1,1,3,1,1,3,1,1,3,1,1,3
???.###????.###????.###????.###????.### ||| [1 1 3 1 1 3 1 1 3 1 1 3 1 1 3]