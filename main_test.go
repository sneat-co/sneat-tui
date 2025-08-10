package main

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

type fakeProgram struct{ run func() (tea.Model, error) }

func (f fakeProgram) Run() (tea.Model, error) { return f.run() }

func TestNewProgram_NotNil(t *testing.T) {
	p := newProgram()
	if p == nil {
		t.Fatalf("newProgram() returned nil")
	}
}

func TestMain_RunSuccess_NoExit(t *testing.T) {
	// Save and restore hooks
	oldGet := getProgram
	oldExit := exit
	defer func() {
		getProgram = oldGet
		exit = oldExit
	}()

	exitCalled := false
	exit = func(code int) { exitCalled = true }

	getProgram = func() program {
		return fakeProgram{run: func() (tea.Model, error) { return nil, nil }}
	}

	main()
	if exitCalled {
		t.Fatalf("exit was called on success path, want not called")
	}
}

func TestMain_RunError_ExitCalled(t *testing.T) {
	oldGet := getProgram
	oldExit := exit
	defer func() {
		getProgram = oldGet
		exit = oldExit
	}()

	exitCalled := false
	exitCode := 0
	exit = func(code int) { exitCalled = true; exitCode = code }

	getProgram = func() program {
		return fakeProgram{run: func() (tea.Model, error) { return nil, errors.New("boom") }}
	}

	main()
	if !exitCalled {
		t.Fatalf("exit was not called on error path")
	}
	if exitCode != 1 {
		t.Fatalf("exit called with code %d, want 1", exitCode)
	}
}
