// redsauce
// 1-dimensional cellular au-tomat-a 🍅
// Cass Smith, October 2019

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
	"os/exec"
	"io"
)

var initState = flag.String("state", "", "Initial worldstate. Characters used to represent living and dead cells must match those specified by the alive and dead options.")
var random = flag.Int("random", 0, "Generate a random initial state of the given size.")
var generations = flag.Int("gen", 10, "The number of generations to iterate.")
var quiet = flag.Bool("quiet", false, "Only print the final worldstate.")
var endState = flag.Bool("end", false, "The logical state of cells outside the world. Ignored if wrap is enabled.")
var wrap = flag.Bool("wrap", false, "If true, the ends of the world are connected.")
var alive = flag.String("alive", "1", "Single character used to represent a living cell.")
var dead = flag.String("dead", "0", "Single character used to represent a dead cell.")
var wolfram = flag.Int("wolfram", -1, "The Wolfram code for a 1-D cellular automation rule.")
var ruleStr = flag.String("rule", "", "Rule described as a logical expression.")

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

func boolify(state string, a, d string) []bool {
	var bools []bool
	for _, cell := range state {
		if string(cell) != d && string(cell) != a {
			continue
		}
		bools = append(bools, string(cell) == a)
	}
	return bools
}

func step(state []bool, r map[int]int, w, e bool, b int) []bool {
	var newState []bool
	for i, _ := range state {
		newState = append(newState, map[int]bool{
			1:  true,
			//-1: false,
			0:	false}[r[toInt(getSubState(state, i, w, e, b))]])
			//0:  cell
		//  1: Cell comes to life
		// 	0: Cell dies
		// //	0: Cell's state doesn't change
	}
	return newState
}

func getCellWrap(state []bool, index int, w, e bool) bool {
	if index >= len(state) {
		if w {
			for index >= len(state) {
				index -= len(state)
			}
			return state[index]
		} else {
			return e
		}
	} else if index < 0 {
		if w {
			for index < 0 {
				index += len(state)
			}
			return state[index]
		} else {
			return e
		}
	} else {
		return state[index]
	}
}

func getSubState(state []bool, index int, w, e bool, b int) []bool {
	var subState []bool
	subState = append(subState, getCellWrap(state, index+1, w, e))
	subState = append(subState, getCellWrap(state, index, w, e))
	subState = append(subState, getCellWrap(state, index-1, w, e))
	for i := 2; len(subState) < b; i += 1 {
		subState = append(subState, getCellWrap(state, index+i, w, e))
		subState = append(subState, getCellWrap(state, index-i, w, e))
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

func stringify(state []bool, a, d string) string {
	outString := ""
	for _, cell := range state {
		outString += map[bool]string{true: a, false: d}[cell]
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
			ruleDef[i] = 0
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

func getParsedRule(ruleString string) (map[int]int, int) {
	parser := exec.Command("./ruledef")
	parserIn, err := parser.StdinPipe()
	if err != nil {
		die(err.Error())
	}
	go func() {
		defer parserIn.Close()
		io.WriteString(parserIn, ruleString)
	}()
	out, err := parser.CombinedOutput()
	if err != nil {
		die(err.Error())
	}
	split := strings.Split(string(out), "\n")
	ruleWidth := int(math.Log2(float64(len(split))))
	if ruleWidth % 2 == 0 {
		ruleWidth += 1
	}
	ruleDef := make(map[int]int)
	for i, v := range split {
		if v == "1" || v == "0" {
			ruleDef[i] = map[string]int{"1": 1, "0": 0}[v]
		}
	}
	return ruleDef, ruleWidth
}

func main() {
	flag.Parse()
	validateRepresentations(*alive, *dead)
	var rule map[int]int
	var bitWidth int
	if *wolfram >= 0 {
		rule = unpackWolfram(*wolfram)
		bitWidth = 3
	} else if len(*ruleStr) > 0 {
		rule, bitWidth = getParsedRule(*ruleStr)
	} else {
		rule = unpackWolfram(110)
		bitWidth = 3
	}
	var worldState []bool
	if *random > 0 {
		worldState = randState(*random)
	} else if len(*initState) > 0 {
		worldState = boolify(*initState, *alive, *dead)
	} else {
		die("No initial state specified.")
	}
	validate(worldState, *generations)
	if !*quiet {
		fmt.Println(stringify(worldState, *alive, *dead))
	}
	for i := 0; i < *generations; i++ {
		worldState = step(worldState, rule, *wrap, *endState, bitWidth)
		if !*quiet || i == *generations-1 {
			fmt.Println(stringify(worldState, *alive, *dead))
		}
	}
}
