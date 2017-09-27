package fsm

import (
	"testing"
)

func TestTransition(t *testing.T) {
	fsm := NewFSM()

	fsm.AddTransition("sharing", "liking")
	fsm.AddTransition("sharing", "uploading")
	fsm.AddTransition("liking", "discovering")
	fsm.AddTransition("discovering", "sharing")

	err := fsm.Start("sharing")
	if err != nil {
		t.Fatal(err)
	}

	// test simple transitioning
	err = fsm.Transition("liking")
	if err != nil {
		t.Fatal("got err when transitioning, expected no error")
	}

	err = fsm.Transition("discovering")
	if err != nil {
		t.Fatal("got err when transitioning, expected no error")
	}

	// test transitioning to multiple states
	err = fsm.Transition("sharing")
	if err != nil {
		t.Fatal("got err when transitioning, expected no error")
	}

	err = fsm.Transition("uploading")
	if err != nil {
		t.Fatal("got err when transitioning, expected no error")
	}

	err = fsm.Transition("invalid state")
	if err == nil {
		t.Fatal("expected err when transitioning, got nil")
	}
}

func TestFuncEnter(t *testing.T) {
	fsm := NewFSM()

	sharingEnterCalled := false
	sharingExitedCalled := false

	likingEnterCalled := false
	likingExitedCalled := false

	fsm.AddTransition("sharing", "liking")
	fsm.AddTransition("liking", "viewing")

	fsm.OnEnter("sharing", func() error {
		if sharingExitedCalled {
			t.Fatal("exit called before enter")
		}
		sharingEnterCalled = true
		return nil
	})

	fsm.OnExit("sharing", func() error {
		if !sharingEnterCalled {
			t.Fatal("exit called before enter")
		}

		sharingExitedCalled = true
		return nil
	})

	fsm.OnEnter("liking", func() error {
		if likingExitedCalled {
			t.Fatal("exit called before enter")
		}

		likingEnterCalled = true
		return nil
	})

	fsm.OnExit("liking", func() error {
		if !likingEnterCalled {
			t.Fatal("exit called before enter")
		}

		likingExitedCalled = true
		return nil
	})

	err := fsm.Start("sharing")
	if err != nil {
		t.Fatal(err)
	}

	err = fsm.Transition("liking")
	if err != nil {
		t.Fatal(err)
	}

	err = fsm.Transition("viewing")
	if err != nil {
		t.Fatal(err)
	}

	if !sharingEnterCalled || !sharingExitedCalled || !likingEnterCalled || !likingExitedCalled {
		t.Fatal("sharing and liking funcs not called correctly")
	}

}
