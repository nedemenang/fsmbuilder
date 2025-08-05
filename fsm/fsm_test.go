package fsm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFSMBuilder(t *testing.T) {
	tests := []struct {
		name        string
		setupFSM    func() (*Builder, error)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid FSM Build",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("S0", "S1", "S2").
					AddSymbols("0", "1").
					SetInitialState("S0")
				if err != nil {
					return nil, err
				}
				builder, err = builder.AddFinalStates("S0", "S1", "S2")
				if err != nil {
					return nil, err
				}

				transitions := []map[TransitionKey]State{
					{{State: "S0", Symbol: "0"}: "S0"},
					{{State: "S0", Symbol: "1"}: "S1"},
					{{State: "S1", Symbol: "0"}: "S2"},
					{{State: "S1", Symbol: "1"}: "S0"},
					{{State: "S2", Symbol: "0"}: "S1"},
					{{State: "S2", Symbol: "1"}: "S2"},
				}

				builder, err = builder.AddTransitions(transitions)
				if err != nil {
					return nil, err
				}

				_, err = builder.Build()
				if err != nil {
					return nil, err
				}
				return builder, nil
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Error on Invalid Initial State",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").
					SetInitialState("S0")
				if err != nil {
					return nil, err
				}
				return builder, err
			},
			expectError: true,
			errorMsg:    "state S0 not in state set",
		},
		{
			name: "Error on Invalid final States",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").
					SetInitialState("S0")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddFinalStates("S0", "S1")
				if err != nil {
					return nil, err
				}
				return builder, err
			},
			expectError: true,
			errorMsg:    "state S0 not in state set",
		},
		{
			name: "Error on Invalid Transitions, missing state",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").
					SetInitialState("q0")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddFinalStates("q0", "q1")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddTransition("S0", "0", "q1")
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "state S0 not in state set",
		},
		{
			name: "Error on Invalid Transitions, missing next state",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").
					SetInitialState("q0")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddFinalStates("q0", "q1")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddTransition("q0", "0", "s1")
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "next state s1 not in state set",
		},
		{
			name: "Error on Invalid Transitions, symbol not in alphabet",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").
					SetInitialState("q0")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddFinalStates("q0", "q1")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddTransition("q0", "4", "q1")
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "symbol 4 not in alphabet",
		},
		{
			name: "Error on Invalid Transitions, transition already exists",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").
					SetInitialState("q0")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddFinalStates("q0", "q1")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddTransition("q0", "0", "q1")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddTransition("q0", "0", "q1")
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "transition δ(q0, 0) already defined",
		},
		{
			name: "Error on build, FSM must have at least one state",
			setupFSM: func() (*Builder, error) {
				builder := NewBuilder()

				_, err := builder.Build()
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "FSM must have at least one state",
		},
		{
			name: "Error on build, FSM must have at least one symbol in alphabet",
			setupFSM: func() (*Builder, error) {
				builder := NewBuilder().AddStates("q0", "q1", "q2")

				_, err := builder.Build()
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "FSM must have at least one symbol in alphabet",
		},
		{
			name: "Error on build, FSM must have an initial state",
			setupFSM: func() (*Builder, error) {
				builder := NewBuilder().AddStates("q0", "q1", "q2").
					AddSymbols("0", "1")

				_, err := builder.Build()
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "FSM must have an initial state",
		},
		{
			name: "Error on build, FSM must have at least one final state",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").SetInitialState("q0")
				if err != nil {
					return nil, err
				}

				_, err = builder.Build()
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "FSM must have at least one final state",
		},
		{
			name: "Error on build, Initial state must be in state set",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().AddStates("q0", "q1", "q2").
					AddSymbols("0", "1").AddFinalStates("q0", "q1", "q2")

				builder.fsm.initialState = "m0"
				_, err = builder.Build()
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "initial state must be in state set",
		},
		{
			name: "Error on build, a transition is not defined",
			setupFSM: func() (*Builder, error) {
				builder, err := NewBuilder().
					AddStates("q0").
					AddSymbols("0", "1").
					SetInitialState("q0")

				builder, err = builder.AddFinalStates("q0")
				if err != nil {
					return nil, err
				}

				builder, err = builder.AddTransition("q0", "0", "q0")
				if err != nil {
					return nil, err
				}

				_, err = builder.Build()
				if err != nil {
					return nil, err
				}

				return builder, err
			},
			expectError: true,
			errorMsg:    "transition δ(q0, 1) is not defined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.setupFSM()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestFSMOperations(t *testing.T) {
	builder, _ := NewBuilder().
		AddStates("q0", "q1").
		AddSymbols("a", "b").
		SetInitialState("q0")

	builder, _ = builder.AddFinalStates("q1")
	var transitions = []map[TransitionKey]State{
		{{State: "q0", Symbol: "a"}: "q1"},
		{{State: "q0", Symbol: "b"}: "q0"},
		{{State: "q1", Symbol: "a"}: "q1"},
		{{State: "q1", Symbol: "b"}: "q0"},
	}
	builder, _ = builder.AddTransitions(transitions)
	fsm, _ := builder.Build()

	t.Run("Step operations", func(t *testing.T) {
		fsm.Reset()
		if fsm.CurrentState() != "q0" {
			t.Errorf("Expected initial state q0, got %s", fsm.CurrentState())
		}

		err := fsm.step("a")
		if err != nil {
			t.Errorf("Unexpected error on step: %v", err)
		}
		if fsm.CurrentState() != "q1" {
			t.Errorf("Expected state q1 after input 'a', got %s", fsm.CurrentState())
		}
	})

	t.Run("ProcessInput operations", func(t *testing.T) {
		fsm.Reset()
		if fsm.CurrentState() != "q0" {
			t.Errorf("Expected initial state q0, got %s", fsm.CurrentState())
		}

		isFinal, _ := fsm.ProcessInput("ab")
		if isFinal {
			t.Error("Expected not final state after processing 'ab'")
		}
		if fsm.CurrentState() != "q0" {
			t.Errorf("Expected state q0 after input 'ab', got %s", fsm.CurrentState())
		}

		isFinal, _ = fsm.ProcessInput("aa")
		if !isFinal {
			t.Error("Expected final state after processing 'aa'")
		}
		if fsm.CurrentState() != "q1" {
			t.Errorf("Expected state q1 after input 'aa', got %s", fsm.CurrentState())
		}
	})

	t.Run("Invalid symbol in processInput", func(t *testing.T) {
		fsm.Reset()
		_, err := fsm.ProcessInput("c")
		if err == nil {
			t.Error("Expected error for invalid symbol 'c', got none")
		} else if err.Error() != "symbol c not in alphabet" {
			t.Errorf("Expected error message 'symbol c not in alphabet', got %v", err)
		}
	})

	t.Run("Get States, alphabet and final states", func(t *testing.T) {
		expectedStates := []State{"q0", "q1"}
		expectedAlphabet := []Symbol{"a", "b"}
		expectedFinalStates := []State{"q1"}

		if !reflect.DeepEqual(fsm.GetStates(), expectedStates) {
			t.Errorf("Expected states %v, got %v", expectedStates, fsm.GetStates())
		}
		if !reflect.DeepEqual(fsm.GetAlphabet(), expectedAlphabet) {
			t.Errorf("Expected alphabet %v, got %v", expectedAlphabet, fsm.GetAlphabet())
		}
		if !reflect.DeepEqual(fsm.GetFinalStates(), expectedFinalStates) {
			t.Errorf("Expected final states %v, got %v", expectedFinalStates, fsm.GetFinalStates())
		}
	})
}

func NewModThreeFSM() (*FSM, error) {
	builder, err := NewBuilder().
		AddStates("s0", "s1", "s2").
		AddSymbols("0", "1").
		SetInitialState("s0")
	if err != nil {
		return nil, fmt.Errorf("error setting initial state: %w", err)
	}

	builder, err = builder.AddFinalStates("s0", "s1", "s2")
	if err != nil {
		return nil, fmt.Errorf("error adding final states: %w", err)
	}
	var transitions = []map[TransitionKey]State{
		{{State: "s0", Symbol: "0"}: "s0"},
		{{State: "s0", Symbol: "1"}: "s1"},
		{{State: "s1", Symbol: "0"}: "s2"},
		{{State: "s1", Symbol: "1"}: "s0"},
		{{State: "s2", Symbol: "0"}: "s1"},
		{{State: "s2", Symbol: "1"}: "s2"},
	}

	builder, err = builder.AddTransitions(transitions)
	if err != nil {
		return nil, fmt.Errorf("error adding transitions: %w", err)
	}

	ModThreeAutomaton, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("error building FSM: %w", err)
	}

	return ModThreeAutomaton, nil
}

func TestModThreeFunction(t *testing.T) {
	modThreeOperation := func(input string) (int, error) {
		ModThreeAutomaton, _ := NewModThreeFSM()
		_, err := ModThreeAutomaton.ProcessInput(input)
		if err != nil {
			return 0, fmt.Errorf("error processing output: %w", err)
		}

		switch ModThreeAutomaton.CurrentState() {
		case "s0":
			return 0, nil // 0 % 3 = 0
		case "s1":
			return 1, nil // 1 % 3 = 1
		case "s2":
			return 2, nil // 2 % 3 = 2
		default:
			return 0, fmt.Errorf("unknown state: %s", ModThreeAutomaton.CurrentState())
		}

	}

	testCases := []struct {
		binary   string
		expected int
	}{
		{"0", 0},    // 0 % 3 = 0
		{"1", 1},    // 1 % 3 = 1
		{"10", 2},   // 2 % 3 = 2
		{"11", 0},   // 3 % 3 = 0
		{"100", 1},  // 4 % 3 = 1
		{"101", 2},  // 5 % 3 = 2
		{"110", 0},  // 6 % 3 = 0
		{"111", 1},  // 7 % 3 = 1
		{"1000", 2}, // 8 % 3 = 2
		{"1001", 0}, // 9 % 3 = 0
		{"1010", 1}, // 10 % 3 = 1
		{"1011", 2}, // 11 % 3 = 2
		{"1100", 0}, // 12 % 3 = 0
		{"1101", 1}, // 13 % 3 = 1
		{"1110", 2}, // 14 % 3 = 2
		{"1111", 0}, // 15 % 3 = 0
	}

	for _, tc := range testCases {
		answer, err := modThreeOperation(tc.binary)
		if err != nil {
			t.Errorf("ModThree(%s) failed: %v", tc.binary, err)
		}
		if answer != tc.expected {
			t.Errorf("ModThree(%s) = %d, expected %d", tc.binary, answer, tc.expected)
		}
	}
}

func BenchmarkModThreeFSM(b *testing.B) {
	fsm, _ := NewModThreeFSM()
	input := "1101010111001010"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fsm.Reset()
		fsm.ProcessString(input)
	}
}

func BenchmarkModThreeFunction(b *testing.B) {
	fsm, _ := NewModThreeFSM()
	input := "1101010111001010"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fsm.Reset()
		fsm.ProcessString(input)
	}
}
