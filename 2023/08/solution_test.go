package main

import (
	"testing"
)

func TestSolution2(t *testing.T) {
	result, err := getResultPart1("test_data2.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	expected := 2
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestSolution6(t *testing.T) {
	result, err := getResultPart1("test_data6.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	expected := 6
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// func TestSolutionPart2(t *testing.T) {
// 	result, err := getResultPart2("test_data_part2.txt")
// 	if err != nil {
// 		t.Errorf("Function failed %s", err)
// 	}
// 	expected := 6
// 	if result != expected {
// 		t.Errorf("Expected %d, got %d", expected, result)
// 	}
// }

type TestData struct {
	startNode           string
	instructions        string
	expectedCycleOffset int
	expectedCycleLength int
	expectedCycleNode   string
	expectedCycleIndex  int
}

// func TestFindCircles(t *testing.T) {
// 	directions := map[string][2]string{
// 		"AAA": [2]string{"AAB", "AAC"},
// 		"AAB": [2]string{"AAX", "AAD"},
// 		"AAC": [2]string{"AAB", "AAC"},
// 		"AAX": [2]string{"AAZ", "AAZ"},
// 		"AAZ": [2]string{"AAX", "AAX"},
// 		"AA0": [2]string{"AA0", "AA0"},
// 		"AAD": [2]string{"AA0", "AAA"},
// 	}
// 	testCases := []TestData{
// 		{startNode: "AAA", instructions: "L", expectedCycleOffset: 2, expectedCycleLength: 2, expectedCycleNode: "AAX", expectedCycleIndex: 0},
// 		{startNode: "AAB", instructions: "L", expectedCycleOffset: 1, expectedCycleLength: 2, expectedCycleNode: "AAX", expectedCycleIndex: 0},
// 		{startNode: "AAA", instructions: "LL", expectedCycleOffset: 2, expectedCycleLength: 2, expectedCycleNode: "AAX", expectedCycleIndex: 0},
// 		{startNode: "AAA", instructions: "R", expectedCycleOffset: 1, expectedCycleLength: 1, expectedCycleNode: "AAC", expectedCycleIndex: 0},
// 	}
// 	for i, testCase := range testCases {

// 		cycle, err := findCycle(testCase.startNode, testCase.instructions, directions)

// 		if err != nil {
// 			t.Errorf("Function failed %s", err)
// 		}
// 		expected := Cycle{
// 			offset:           testCase.expectedCycleOffset,
// 			length:           testCase.expectedCycleLength,
// 			node:             testCase.expectedCycleNode,
// 			instructionIndex: testCase.expectedCycleIndex,
// 		}
// 		if cycle != expected {
// 			t.Errorf("Failure at %d: Expected %#v, got %#v", i, expected, cycle)
// 		}
// 	}

// }

type GCDTestCase struct {
	a        int
	b        int
	expected int
}

func TestGCD(t *testing.T) {
	testCases := []GCDTestCase{
		{5, 5, 5},
		{5, 10, 5},
		{0, 2, 2},
		{81, 18, 9},
	}
	for _, testCase := range testCases {
		result := gcd(testCase.a, testCase.b)
		if result != testCase.expected {
			t.Errorf("GDC(%d, %d) expected %d, got %d", testCase.a, testCase.b, testCase.expected, result)
		}
	}
}

func TestLCM(t *testing.T) {
	testCases := []GCDTestCase{
		{5, 5, 5},
		{5, 10, 10},
		{3, 5, 15},
		{14, 18, 7 * 2 * 9},
	}
	for _, testCase := range testCases {
		result := lcm(testCase.a, testCase.b)
		if result != testCase.expected {
			t.Errorf("LCM(%d, %d) expected %d, got %d", testCase.a, testCase.b, testCase.expected, result)
		}
	}
}
