package sneatui

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
)

func TestUnsigned_EnterOnAbout_OpensAboutViaApp(t *testing.T) {
	m := InitialModel().(appModel)
	// Move selection down to About
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	am := model.(appModel)
	// Enter
	model, cmd := am.Update(tea.KeyMsg{Type: tea.KeyEnter})
	am = model.(appModel)
	if cmd != nil {
		msg := cmd()
		model, _ = am.Update(msg)
		am = model.(appModel)
	}
	if am.active != "about" {
		t.Fatalf("after Enter on About, active=%q, want about", am.active)
	}
}

func TestAbout_EscReturnsToMenu_ViaApp(t *testing.T) {
	m := InitialModel().(appModel)
	// Go to About
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	am := model.(appModel)
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
	if am.active != "unsigned" {
		t.Fatalf("after ESC from about, active=%q, want unsigned", am.active)
	}
}

func TestAbout_ViewContainsContent(t *testing.T) {
	a := newAboutModel().(aboutModel)
	// trigger Init to load content
	cmd := a.Init()
	if cmd != nil {
		msg := cmd()
		// deliver the loaded message to update state
		model, _ := a.Update(msg)
		a = model.(aboutModel)
	}
	view := a.View()
	if !strings.Contains(view, "Sneat") {
		t.Fatalf("about view does not contain expected content, got:\n%s", view)
	}
}

func TestUnsigned_SelectionPreserved_AfterAboutEsc(t *testing.T) {
	m := InitialModel().(appModel)
	// Move selection to About
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	am := model.(appModel)
	idxBefore := am.unsigned.list.Index()
	// Enter to open About
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
