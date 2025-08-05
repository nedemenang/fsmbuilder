# fsmbuilder

`fsmbuilder` is a Go library for building and running finite state machines (FSMs).

## Features

- Builder pattern for FSM construction
- Support for custom states, symbols, and transitions
- Input processing to determine final state acceptance


Below is an example of how to define and use a finite state machine with `fsmbuilder`:

```go
package main

import (
    "fmt"
    "github.com/username/fsmbuilder/fsm"
)

func main() {
    // Create a new FSM builder
    b := fsm.NewBuilder().
        AddStates("q0", "q1").
        AddSymbols("a", "b")

    // Set initial and final states
    _, _ = b.SetInitialState("q0")
    _, _ = b.AddFinalStates("q1")

    // Define transitions
    transitions := []map[fsm.TransitionKey]fsm.State{
        {{State: "q0", Symbol: "a"}: "q1"},
        {{State: "q0", Symbol: "b"}: "q0"},
        {{State: "q1", Symbol: "a"}: "q1"},
        {{State: "q1", Symbol: "b"}: "q0"},
    }
    _, _ = b.AddTransitions(transitions)

    // Build the FSM
    f, err := b.Build()
    if err != nil {
        panic(err)
    }

    // Test input strings
    inputs := []string{"a", "b", "aa", "ab", "ba", "bb"}
    for _, input := range inputs {
        accepted := f.ProcessOutput(input)
        fmt.Printf("Input: %s, Accepted: %v\n", input, accepted)
    }
}
```

The main.go file includes an implementation of a Modulo Three FSM with error handling and input processing.