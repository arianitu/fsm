package fsm

import (
	"errors"
)

type transitionFunc func() error

// FSM is an implementation of a finite state machine.
type FSM struct {
	// currentState is the current state that the FSM is in
	currentState string

	// transitionMap maps a single state to N valid transition states
	transitionMap map[string][]string

	// enterFuncs is a map to the function that will be called when a state is entered
	enterFuncs map[string]transitionFunc

	// exitFuncs is a map to the function that will be called when a state is exited
	exitFuncs map[string]transitionFunc
}

// NewFSM creates a new finite state machine. The starting state is nil and you must
// start the FSM by calling .Start() after adding all transitions and the proper entrance functions.
func NewFSM() *FSM {
	return &FSM{
		currentState:  "",
		transitionMap: make(map[string][]string, 0),
		enterFuncs:    make(map[string]transitionFunc, 0),
		exitFuncs:     make(map[string]transitionFunc, 0),
	}
}

// AddTransition adds a transition from stateA to stateB. States that transition to themselves
// are valid.
func (f *FSM) AddTransition(stateA string, stateB string) {
	_, ok := f.transitionMap[stateA]
	if !ok {
		f.transitionMap[stateA] = make([]string, 0)
	}

	// do nothing if transition already exists
	for _, transition := range f.transitionMap[stateA] {
		if transition == stateB {
			return
		}
	}

	f.transitionMap[stateA] = append(f.transitionMap[stateA], stateB)

}

// OnEnter adds a function that should be called when we enter a transition
func (f *FSM) OnEnter(state string, fun transitionFunc) error {
	f.enterFuncs[state] = fun
	return nil
}

// OnExit adds a function that should be called when we exit a transition
func (f *FSM) OnExit(state string, fun transitionFunc) error {
	f.exitFuncs[state] = fun
	return nil
}

// Start initializes the finite state machine with the start state. The start state
// onEnter func will be called.
func (f *FSM) Start(startState string) error {
	fun, ok := f.enterFuncs[startState]
	if ok {
		err := fun()
		if err != nil {
			return err
		}
	}

	f.currentState = startState
	return nil
}

// Transition transitions from the current FSM state to nextState. If the transition is not
// a valid transition, error will be non nil. The onExit func of the current state is called
// and the onEnter func of the next state is called
func (f *FSM) Transition(nextState string) error {
	transitions, ok := f.transitionMap[f.currentState]
	if !ok {
		return errors.New("Could not transition to " + nextState + ". State does not exist")
	}

	found := false
	for _, transition := range transitions {
		if transition == nextState {
			found = true
			break
		}
	}

	if !found {
		return errors.New("Invalid transition from " + f.currentState + " to " + nextState)
	}

	// call exit funcs
	exitFun, ok := f.exitFuncs[f.currentState]
	if ok {
		err := exitFun()
		if err != nil {
			return err
		}
	}

	// call enter funcs
	enterFun, ok := f.enterFuncs[nextState]
	if ok {
		err := enterFun()
		if err != nil {
			return err
		}
	}
	f.currentState = nextState

	return nil
}
