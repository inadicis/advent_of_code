package main

import "testing"

/* TEST DATA

seeds: 79 14 55 13
seed-to-soil map:
50 98 2				// destination range start, source range start, range length, so 98 -> 50 and 99 -> 51
52 50 48			// 50..98 -> 52..100
					// and as there is no map, 1..50 -> 1..50 (same number)

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
*/

/*
Notes
seed -> soil -> fertilizer -> water -> light -> temperatur -> humidity -> location
could simplify each map with a directed graph (probably even a tree if for source a != b, destinations m[a] != m[b])
-----
that gets harder with ranges

maps are in order for the chain, we can simply process one after each other

*/

func TestSolution(t *testing.T) {
	result := solution5("test_data.txt")
	if result != 35 {
		t.Errorf("Expected %d, got %d", 35, result)
	}
}
