package main

import (
	"testing"
)

func TestSolution(t *testing.T) {
	result, err := amountWaysToBeatWR("test_data.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	expected := 71503
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
