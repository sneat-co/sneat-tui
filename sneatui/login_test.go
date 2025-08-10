package sneatui

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
)

func TestUnsigned_EnterOnSignIn_OpensLoginViaApp(t *testing.T) {
	m := InitialModel().(appModel)
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am := model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.active != "login" {
		t.Fatalf("after Enter on Sign-in, active=%q, want login", am.active)
	}
}

func TestLogin_EscReturnsToUnsignedMenu_ViaApp(t *testing.T) {
	m := InitialModel().(appModel)
	// Go to login
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am := model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	// Send ESC from login
	model, cmd = am.Update(tea.KeyMsg{Type: tea.KeyEsc})
	am = model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.active != "unsigned" {
		t.Fatalf("after ESC from login, active=%q, want unsigned", am.active)
	}
}

func TestLogin_ViewNotEmpty(t *testing.T) {
	m := newLoginModel().(loginModel)
	if m.View() == "" {
		t.Fatalf("login view is empty")
	}
}

func TestUnsigned_SelectionPreserved_AfterLoginEsc(t *testing.T) {
	am := InitialModel().(appModel)
	// Move selection down
	model, _ := am.Update(tea.KeyMsg{Type: tea.KeyDown})
	am = model.(appModel)
	idxBefore := am.unsigned.list.Index()
	// Enter to open login
	model, cmd := am.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am = model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	// ESC back
	model, cmd = am.Update(tea.KeyMsg{Type: tea.KeyEsc})
	am = model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.unsigned.list.Index() != idxBefore {
		t.Fatalf("selection index after return = %d, want %d", am.unsigned.list.Index(), idxBefore)
	}
	view := am.View()
	if !strings.Contains(view, "Sneat.app - main menu") || !strings.Contains(view, "Sign-in") {
		t.Fatalf("menu view after return missing expected content; view=\n%s", view)
	}
}
