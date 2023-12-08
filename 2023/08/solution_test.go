package main

import (
	"testing"
)

func TestSolution2(t *testing.T) {
	result, err := getResult("test_data2.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	expected := 2
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestSolution6(t *testing.T) {
	result, err := getResult("test_data6.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	expected := 6
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
