package sneatui

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
)

func TestLogin_EnterOnPassword_NavigatesToSignedMenu(t *testing.T) {
	m := InitialModel().(appModel)
	// Enter on Sign-in to go to login
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am := model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.active != "login" {
		t.Fatalf("expected to be on login, got %q", am.active)
	}
	// Press Enter twice: first moves to password, second submits
	model, _ = am.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am = model.(appModel)
	model, cmd = am.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am = model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.active != "signed" {
		t.Fatalf("after login submit, active=%q, want signed", am.active)
	}
	if !strings.Contains(am.View(), "Sign-out") {
		t.Fatalf("signed-in menu view does not contain Sign-out")
	}
}

func TestSignedIn_SignOut_ReturnsUnsigned(t *testing.T) {
	// Move to signed-in first
	m := InitialModel().(appModel)
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am := model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	// Enter twice to submit login
	model, _ = am.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am = model.(appModel)
	model, cmd = am.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am = model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.active != "signed" {
		t.Fatalf("expected signed before sign-out, got %q", am.active)
	}
	// Now select Sign-out (default index 0 is Calendar; move cursor to last index)
	// We'll repeatedly send down keys to reach Sign-out
	for i := 0; i < 3; i++ { // 3 moves from 0 -> 3
		model, _ = am.Update(tea.KeyMsg{Type: tea.KeyDown})
		am = model.(appModel)
	}
	model, cmd = am.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am = model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.active != "unsigned" {
		t.Fatalf("after Sign-out, active=%q, want unsigned", am.active)
	}
}
