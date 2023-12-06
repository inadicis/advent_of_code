package main

import (
	"reflect"
	"testing"
)

func TestSolution(t *testing.T) {
	result, err := amountWaysToBeatWR("test_data.txt")
	if err != nil {
		t.Errorf("Function failed %s", err)
	}
	expected := []int{4, 8, 9}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
