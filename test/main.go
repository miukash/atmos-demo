package main

import "fmt"

func main() {
	nfa, err := Compile("x(ab)y")
	if err != nil {
		panic(err)
	}

	input := "xaby"

	fmt.Println("Input: ", input)
	fmt.Println("Match: ", nfa.Match(input))
	fmt.Println("Masked:", nfa.Mask(input))

	result, ok := nfa.MatchCapture(input)
	if ok {
		fmt.Printf("Capture: [%d,%d)\n", result.Start, result.End)
	}
}
