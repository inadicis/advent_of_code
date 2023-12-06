package main

import "testing"

func TestSolution(t *testing.T) {
	result, err := amountWays("test_data.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	expected :=  []int{4, 8, 9}
	if result != {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
