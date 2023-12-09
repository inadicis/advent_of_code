package main

import (
	"reflect"
	"testing"
)

func TestExtractData(t *testing.T) {
	fileName := "test_data.txt"
	expected := [][]int{
		{0, 3, 6, 9, 12, 15},
		{1, 3, 6, 10, 15, 21},
		{1, 0, 13, 16, 21, 30, 45},
		{55, -95, -349, -669, -918, -793, 264, 3180},
	}
	actual, err := ExtractData(fileName)
	if err != nil {
		t.Errorf("ExtractData failed: %s", err)
	}
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

type DeepTestCase struct {
	row []int 
	deepness int 
	found bool
}

func TestFindDeepnessOf0Row(t *testing.T) {
	testCases := []DeepTestCase{
		{row: []int{0, 3, 6, 9, 12, 15}, deepness: 2, found: true},
		{row: []int{1, 3, 6, 10, 15, 21}, deepness: 3, found: true},
		{row: []int{10  13  16  21  30  45  68}, deepness: 3, found: true},
		{row: []int{1, 2}, deepness: 0, found: false},
	}
	for i, expected := range testCases {
		d, found := FindDeepnessOf0Row(expected.row)
		if found != expected.found {
			t.Errorf("case %d: expected found %v, got %v", i, expected.found, found)
		}
		if expected.found && expected.deepness != d {
			t.Errorf("case %d: expected found ad deepness %d, got %d", i, expected.deepness, d)

		}
	}
}
