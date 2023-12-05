package main

import "testing"

func TestSolution(t *testing.T) {
	result := solution5("test_data.txt")
	if result != 35 {
		t.Errorf("Expected %d, got %d", 35, result)
	}
}
