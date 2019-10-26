// redsauce
// 1-dimensional cellular au-tomat-a 🍅
// Cass Smith, October 2019

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

var initState = flag.String("state", "", "Initial worldstate. Use 1s and 0s to specify living and dead cells.")
var generations = flag.Int("gen", 10, "The number of generations to iterate.")

// Rule 30
var rule = map[int]int{
	0: -1,
	1: 1,
	2: 1,
	3: 1,
	4: 1,
	5: -1,
	6: -1,
	7: -1}

func die(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

func validate(state []bool, gen int) {
	if len(state) < 1 {
		die("World must contain at least one cell.")
	}
	if gen < 0 {
		die("Cannot simulate backward in time.")
	}
	// No failure conditions met, resuming
}

func boolify(state string) []bool {
	var bools []bool
	for _, cell := range state {
		if cell != '0' && cell != '1' {
			continue
		}
		bools = append(bools, cell == '1')
	}
	return bools
}

func step(state []bool, r map[int]int) []bool {
	var newState []bool
	for i, cell := range state {
		newState = append(newState, map[int]bool{
			1:  true,
			-1: false,
			0:  cell}[r[toInt(getSubState(state, i))]])
	}
    return newState
}

func getSubState(state []bool, index int) []bool {
    var subState []bool
    if index <= 0 {
        subState = append(subState, false)
    } else {
        subState = append(subState, state[index-1])
    }
    subState = append(subState, state[index])
    if index + 1 >= len(state) {
        subState = append(subState, false)
    } else {
        subState = append(subState, state[index+1])
    }
	return subState
}

func toInt(subState []bool) int {
	acc := 0
	for i, c := range subState {
		if c {
			acc += int(math.Pow(2.0, float64(i)))
		}
	}
	return acc
}

func stringify(state []bool) string {
    outString := ""
    for _, cell := range state {
        outString += map[bool]string{true: "1", false: "0"}[cell]
    }
    return outString
}

func main() {
	flag.Parse()
	var worldState = boolify(*initState)
	validate(worldState, *generations)
	fmt.Println(stringify(worldState))
    for i := 0; i < *generations; i++ {
        worldState = step(worldState, rule)
        fmt.Println(stringify(worldState))
    }
}
