package main

import (
	"fmt"
	"github.com/nedemenang/fsmbuilder/fsm"
)

func main() {
	builder, err := fsm.NewBuilder().
		AddStates("S0", "S1", "S2").
		AddSymbols("0", "1").
		SetInitialState("S0")
	if err != nil {
		fmt.Println("Error setting initial state:", err)
		return
	}

	builder, err = builder.AddFinalStates("S0", "S1", "S2")
	if err != nil {
		fmt.Println("Error adding final states:", err)
		return
	}
	var transitions = []map[fsm.TransitionKey]fsm.State{
		{{State: "S0", Symbol: "0"}: "S0"},
		{{State: "S0", Symbol: "1"}: "S1"},
		{{State: "S1", Symbol: "0"}: "S2"},
		{{State: "S1", Symbol: "1"}: "S0"},
		{{State: "S2", Symbol: "0"}: "S1"},
		{{State: "S2", Symbol: "1"}: "S2"},
	}

	builder, err = builder.AddTransitions(transitions)
	if err != nil {
		fmt.Println("Error adding transitions:", err)
		return
	}

	moduloThreeAutomaton, err := builder.Build()
	if err != nil {
		fmt.Println("Error building FSM:", err)
		return
	}
	// Example usage of the FSM
	_, err = moduloThreeAutomaton.ProcessInput("110")
	if err != nil {
		fmt.Println("Error processing output:", err)
		return
	}
	switch moduloThreeAutomaton.CurrentState() {
	case "S0":
		fmt.Println("110 => S0 (mod 3 = 0)")
	case "S1":
		fmt.Println("110 => S1 (mod 3 = 1)")
	case "S2":
		fmt.Println("110 => S2 (mod 3 = 2)")
	default:
		fmt.Println("Unknown state:", moduloThreeAutomaton.CurrentState())
	}

	moduloThreeAutomaton.Reset()

	_, err = moduloThreeAutomaton.ProcessInput("1101")
	if err != nil {
		fmt.Println("Error processing output:", err)
		return
	}

	switch moduloThreeAutomaton.CurrentState() {
	case "S0":
		fmt.Println("1101 => S0 (mod 3 = 0)")
	case "S1":
		fmt.Println("1101 => S1 (mod 3 = 1)")
	case "S2":
		fmt.Println("1101 => S2 (mod 3 = 2)")
	default:
		fmt.Println("Unknown state:", moduloThreeAutomaton.CurrentState())
	}
}
