package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	amount, err := amountWaysToBeatWR("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(amount)
	// product := 1
	// for _, amount := range ways {
	// 	product *= amount
	// }
	// fmt.Println(product)
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

func extractInt(line string, skip int) (value int, err error) {
	elements := strings.Split(strings.Trim(line, " \r\n"), " ")
	v := ""
	for i, element := range elements {
		if i < skip || element == "" {
			continue
		}
		v += element
		// rest should be the value
		// values = append(values, v)
	}

	result, err := strconv.Atoi(v)
	return result, err
}

func calcDist(pressedDuration int, raceDuration int) int {
	if pressedDuration <= 0 || pressedDuration >= raceDuration {
		return 0
	}
	return (raceDuration - pressedDuration) * pressedDuration
}

func amountBetterPossibilities(maxTime int, distance int) int {
	// d < (maxTime - x) * x
	// -x^2 + maxTime * x - distance > 0
	delta := maxTime*maxTime - 4*distance
	sqrtDelta := math.Sqrt(float64(delta))
	x0 := int(math.Floor((float64(maxTime)-sqrtDelta)/2 + 1))

	// x1 := (float64(maxTime)+sqrtDelta)/2 + 1
	// fmt.Println("x0 ", x0, "x1", x1)
	return maxTime - 2*(x0-1) - 1

	// for i := 1; i < maxTime; i++ {
	// 	// could do binary search instead of linear search
	// 	if calcDist(i, maxTime) > distance {
	// 		fmt.Println("i ", i)
	// 		return maxTime - 2*(i-1) - 1
	// 	}
	// }
	// return 0
}

func amountWaysToBeatWR(filename string) (amount int, err error) {
	text, err := os.ReadFile(filename)

	if err != nil {
		return 0, err
	}
	// strings.Trim(string(content), " \r\n")
	timesLine, distancesLine, _ := strings.Cut(string(text), "\r\n")
	duration, err := extractInt(timesLine, 1)
	if err != nil {
		return 0, err
	}
	distance, err := extractInt(distancesLine, 1)
	if err != nil {
		return 0, err
	}
	return amountBetterPossibilities(duration, distance), nil

}
