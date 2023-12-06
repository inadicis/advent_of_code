package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	ways, err := amountWaysToBeatWR("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ways)
	product := 1
	for _, amount := range ways {
		product *= amount
	}
	fmt.Println(product)
}

func extractInts(line string, skip int) (values []int, err error) {
	elements := strings.Split(strings.Trim(line, " \r\n"), " ")
	for i, element := range elements {
		if i < skip || element == "" {
			continue
		}
		// rest should be the value
		v, err := strconv.Atoi(element)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func calcDist(pressedDuration int, raceDuration int) int {
	if pressedDuration <= 0 || pressedDuration >= raceDuration {
		return 0
	}
	return (raceDuration - pressedDuration) * pressedDuration
}

func amountBetterPossibilities(time int, distance int) int {
	// x0 = (t - (t^2 + 4*d)^0.5)/2
	// solution = t - 2 * floor((t - (t^2 + 4*d)^0.5)/2 + 1)
	// t := float64(time)
	// d := float64(distance)
	// fmt.Printf("%f - 2*floor((%f - (%f*%f + 4*%f)^0.5)/2 + 1)", t, t, t, t, d)
	// result := int(t - 2*math.Floor((t-math.Sqrt(t*t+4*d))/2+1))
	// fmt.Printf(" t=%d, d=%d, result: %d\n", time, distance, result)
	// return result
	// result := 0
	for i := 1; i < time; i++ {
		if calcDist(i, time) > distance {
			return time - 2*(i-1) - 1
		}
	}
	return 0
	// return result
}

func amountWaysToBeatWR(filename string) (amounts []int, err error) {
	text, err := os.ReadFile(filename)

	if err != nil {
		return []int{}, err
	}
	// strings.Trim(string(content), " \r\n")
	timesLine, distancesLine, _ := strings.Cut(string(text), "\r\n")
	durations, err := extractInts(timesLine, 1)
	if err != nil {
		return nil, err
	}
	distances, err := extractInts(distancesLine, 1)
	if err != nil {
		return nil, err
	}
	fmt.Println(durations, distances)
	for i, duration := range durations {
		distance := distances[i]
		amounts = append(amounts, amountBetterPossibilities(duration, distance))
	}
	return amounts, nil

}

/*
	Solution Ideas:
	- try function for each value v, with 0 < v < duration
	  (as it is symmetric, we can calculate for only the first half then double the amount)
	- mathematical -> analysis of f_n with n is duration ?
	- permutations ? (distance results look like permutatoin amount results:
	 n  0 1   2    3    4   5   6   7
		0 6   10   12   12  10  6   0 (t = 7)
		0 d-1 2d-2 3d-3 3(t-d) - 3
		t*(d - 1)       (d-t)*(t - d)



		6 = 12 - 6  | 12 = 6*2 = t * (d - t + 1)
*/

/*
	With max time t, distance to beat d, we are looking for t_0,
	the minimum pressed diration, so that f(n, t)) > d. Call it g(t, d) -> t_0
	(as t_0 := t/2 is optimal, and distribution is symmetric, to get the amount of ways,
	we can calculate 2 *((t / 2) - t_0) = t - 2*t_0)

	f(n, t) = n * (t-n) for n <= t/2
	n*(t-n) > d  <=> -n^2 + t*n >  d   (always true: t > n)
	-n^2 + t*n - d > 0
	a = -1 , b = t, c = -d
	delta = t^2 + 4*d
	graph is A shaped (as a < 0)
	if delta > 0: there are 2 roots to the polynom
	if delta < 0: unbeatable (impossible) time
	x1 = (t + (t^2 + 4*d)^0.5)/2  (discarded, we want the tiniest)
	x0 = (t - (t^2 + 4*d)^0.5)/2
	solution = t - 2 * (floor((t - (t^2 + 4*d)^0.5)/2 + 1))

	x0 = xmax +- delta?
	f = ax^2 + bx + c
	f' = 2ax + b
	f'(x) = 0 <=> 2ax = -b  <=> x = -b / 2a
	(x - x0)(x - x1) = -x^2 + t*x - d
	x0 = (-b - (delta)^0,5) / 2*a

	x^2 - 4x + 1   | f(-1) = 4; f(0) = 1; f(1) = 0; f(2) = 1
	a: 1 b: -4 c: 1
	delta = 16 - 4 = 12

	x0 = -1 , x1 = 2
	(x + 1)(x - 2) = x^2 -x -2 -> delta = 9 -> root = 3, (-b + 3)/2a = 2
	2x^2 - 2x - 4 -> delta = 4 + 32 = 36 -> 6, (-b + 6)/(2a) = 2


	g(t, d) = floor(d/(t-n) + 1)
*/

// times := strings.Split(timesLine, " ")
// durations := []int{}
// for i, duration := range times {
// 	if i== 0 || duration == " "{ continue } // skip name and empty stinrgs
// 	// rest should be the value
// 	d, err := strconv.Atoi(duration)
// 	if err != nil {
// 		return durations, err
// 	}
// 	durations = append(durations, d)
// }
