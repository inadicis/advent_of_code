package main

import "testing"

func TestSolution(t *testing.T) {
	result, err := solution5("test_data.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	if result != 46 {
		t.Errorf("Expected %d, got %d", 46, result)
	}
}
