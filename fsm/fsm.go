package fsm

import "fmt"

// State represents a state in the finite automaton
type State string

// Symbol represents an input symbol from the alphabet
type Symbol string

// TransitionKey represents a (state, symbol) pair for the transition function
type TransitionKey struct {
	State  State
	Symbol Symbol
}

// FSM represents a finite state machine (Q, Σ, q0, F, δ)
type FSM struct {
	states       map[State]bool          // Q: finite set of states
	alphabet     map[Symbol]bool         // Σ: finite input alphabet
	initialState State                   // q0: initial state
	finalStates  map[State]bool          // F: set of accepting/final states
	transitions  map[TransitionKey]State // δ: transition function Q×Σ→Q
	currentState State                   // current state during execution
}

// Builder provides a fluent interface for building FSMs
type Builder struct {
	fsm *FSM
}

// NewBuilder creates a new FSM builder
func NewBuilder() *Builder {
	return &Builder{
		fsm: &FSM{
			states:      make(map[State]bool),
			alphabet:    make(map[Symbol]bool),
			finalStates: make(map[State]bool),
			transitions: make(map[TransitionKey]State),
		},
	}
}

// AddStates adds multiple states to the FSM
func (b *Builder) AddStates(states ...State) *Builder {
	for _, state := range states {
		b.fsm.states[state] = true
	}
	return b
}

func (b *Builder) AddSymbols(symbols ...Symbol) *Builder {
	for _, symbol := range symbols {
		b.fsm.alphabet[symbol] = true
	}
	return b
}

func (b *Builder) SetInitialState(state State) (*Builder, error) {
	if !b.fsm.states[state] {
		return b, fmt.Errorf("state %s not in state set", state)
	}
	b.fsm.initialState = state
	return b, nil
}

func (b *Builder) AddFinalStates(states ...State) (*Builder, error) {
	for _, state := range states {
		if !b.fsm.states[state] {
			return b, fmt.Errorf("state %s not in state set", state)
		}
		b.fsm.finalStates[state] = true
	}
	return b, nil
}

func (b *Builder) AddTransitions(transitions []map[TransitionKey]State) (*Builder, error) {
	for _, transition := range transitions {
		for key, nextState := range transition {
			b, err := b.AddTransition(key.State, key.Symbol, nextState)
			if err != nil {
				return b, err
			}
		}
	}
	return b, nil
}

func (b *Builder) AddTransition(state State, symbol Symbol, nextState State) (*Builder, error) {
	if !b.fsm.states[state] {
		return b, fmt.Errorf("state %s not in state set", state)
	}
	if !b.fsm.states[nextState] {
		return b, fmt.Errorf("next state %s not in state set", nextState)
	}
	if !b.fsm.alphabet[symbol] {
		return b, fmt.Errorf("symbol %s not in alphabet", symbol)
	}

	key := TransitionKey{State: state, Symbol: symbol}
	if _, exists := b.fsm.transitions[key]; exists {
		return b, fmt.Errorf("transition δ(%s, %s) already defined", state, symbol)
	}

	b.fsm.transitions[key] = nextState
	return b, nil
}

func (b *Builder) Build() (*FSM, error) {
	if len(b.fsm.states) == 0 {
		return nil, fmt.Errorf("FSM must have at least one state")
	}
	if len(b.fsm.alphabet) == 0 {
		return nil, fmt.Errorf("FSM must have at least one symbol in alphabet")
	}
	if b.fsm.initialState == "" {
		return nil, fmt.Errorf("FSM must have an initial state")
	}

	if len(b.fsm.finalStates) == 0 {
		return nil, fmt.Errorf("FSM must have at least one final state")
	}

	if !b.fsm.states[b.fsm.initialState] {
		return nil, fmt.Errorf("initial state must be in state set")
	}

	// Validate that all transitions are defined (total function)
	for state := range b.fsm.states {
		for symbol := range b.fsm.alphabet {
			key := TransitionKey{State: state, Symbol: symbol}
			if _, exists := b.fsm.transitions[key]; !exists {
				return nil, fmt.Errorf("transition δ(%s, %s) is not defined", state, symbol)
			}
		}
	}

	// Initialize current state
	b.fsm.currentState = b.fsm.initialState
	return b.fsm, nil
}

func (f *FSM) Reset() {
	f.currentState = f.initialState
}

func (f *FSM) CurrentState() State {
	return f.currentState
}

func (f *FSM) step(symbol Symbol) error {
	if !f.alphabet[symbol] {
		return fmt.Errorf("symbol %s not in alphabet", symbol)
	}

	key := TransitionKey{State: f.currentState, Symbol: symbol}
	nextState, exists := f.transitions[key]
	if !exists {
		return fmt.Errorf("no transition defined for δ(%s, %s)", f.currentState, symbol)
	}

	f.currentState = nextState
	return nil
}

func (f *FSM) ProcessString(input string) error {
	for _, char := range input {
		symbol := Symbol(char)
		if err := f.step(symbol); err != nil {
			return err
		}
	}
	return nil
}

func (f *FSM) IsFinalState() bool {
	return f.finalStates[f.currentState]
}

func (f *FSM) ProcessInput(input string) (bool, error) {
	f.Reset()
	if err := f.ProcessString(input); err != nil {
		return false, err
	}
	return f.IsFinalState(), nil
}

func (f *FSM) GetStates() []State {
	states := make([]State, 0, len(f.states))
	for state := range f.states {
		states = append(states, state)
	}
	return states
}

// GetAlphabet returns a copy of the alphabet
func (f *FSM) GetAlphabet() []Symbol {
	alphabet := make([]Symbol, 0, len(f.alphabet))
	for symbol := range f.alphabet {
		alphabet = append(alphabet, symbol)
	}
	return alphabet
}

// GetFinalStates returns a copy of the final states set
func (f *FSM) GetFinalStates() []State {
	finalStates := make([]State, 0, len(f.finalStates))
	for state := range f.finalStates {
		finalStates = append(finalStates, state)
	}
	return finalStates
}
