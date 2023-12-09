package main

import (
	"testing"
)

func TestExtractData(t testing.T) {
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
	if expected != actual {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}
