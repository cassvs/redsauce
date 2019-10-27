// redsauce
// 1-dimensional cellular au-tomat-a 🍅
// Cass Smith, October 2019

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"math/rand"
	"time"
)

var initState = flag.String("state", "", "Initial worldstate. Characters used to represent living and dead cells must match those specified by the alive and dead options.")
var random = flag.Int("random", 0, "Generate a random initial state of the given size.")
var generations = flag.Int("gen", 10, "The number of generations to iterate.")
var quiet = flag.Bool("quiet", false, "Only print the final worldstate.")
var endState = flag.Bool("end", false, "The logical state of cells outside the world. Ignored if wrap is enabled.")
var wrap = flag.Bool("wrap", false, "If true, the ends of the world are connected.")
var alive = flag.String("alive", "1", "Single character used to represent a living cell.")
var dead = flag.String("dead", "0", "Single character used to represent a dead cell.")
var wolfram = flag.Int("wolfram", 30, "The Wolfram code for a 1-D cellular automation rule.")

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

func validateRepresentations(aliveRep, deadRep string) {
	if len(aliveRep) != 1 || len(deadRep) != 1 {
		die("Cell representations must be one character.")
	}
}

func boolify(state string) []bool {
	var bools []bool
	for _, cell := range state {
		if string(cell) != *dead && string(cell) != *alive {
			continue
		}
		bools = append(bools, string(cell) == *alive)
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
		//  1: Cell comes to life
		// -1: Cell dies
		// 	0: Cell's state doesn't change
	}
	return newState
}

func getSubState(state []bool, index int) []bool {
	var subState []bool
	if index+1 >= len(state) {
		if *wrap {
			subState = append(subState, state[0])
		} else {
			subState = append(subState, *endState)
		}
	} else {
		subState = append(subState, state[index+1])
	}
	subState = append(subState, state[index])
	if index <= 0 {
		if *wrap {
			subState = append(subState, state[len(state)-1])
			} else {
				subState = append(subState, *endState)
			}
			} else {
				subState = append(subState, state[index-1])
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
		outString += map[bool]string{true: *alive, false: *dead}[cell]
	}
	return outString
}

func unpackWolfram(wolf int) map[int]int {
	if wolf > 255 || wolf < 0 {
		die("Wolfram code must be in the interval [0, 255].")
	}
	ruleDef := make(map[int]int)
	for i := 0; i < 8; i++ {
		if wolf&int(math.Pow(2.0, float64(i))) != 0 {
			ruleDef[i] = 1
		} else {
			ruleDef[i] = -1
		}
	}
	return ruleDef
}

func randState(length int) []bool {
	rand.Seed(time.Now().UnixNano())
	cellStates := []bool{false, true}
	var world []bool
	for i := 0; i < length; i++ {
		world = append(world, cellStates[rand.Intn(len(cellStates))])
	}
	return world
}

func main() {
	flag.Parse()
	validateRepresentations(*alive, *dead)
	rule := unpackWolfram(*wolfram)
	var worldState []bool
	if *random > 0 {
		worldState = randState(*random)
	} else if len(*initState) > 0 {
		worldState = boolify(*initState)
	} else {
		die("No initial state specified.")
	}
	validate(worldState, *generations)
	if !*quiet {
		fmt.Println(stringify(worldState))
	}
	for i := 0; i < *generations; i++ {
		worldState = step(worldState, rule)
		if !*quiet || i == *generations-1 {
			fmt.Println(stringify(worldState))
		}
	}
}
