package sneatui

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestInitialModel_ReturnsMenuUnassigned(t *testing.T) {
	m := InitialModel()
	// Ensure the returned model is our menuUnassigned
	mu, ok := m.(menuUnassigned)
	if !ok {
		// In case the model is a pointer in future, try that too for a clearer error
		if _, okp := m.(*menuUnassigned); okp {
			t.Fatalf("InitialModel returned *menuUnassigned, expected value type menuUnassigned (update tests if implementation changed)")
		}
		t.Fatalf("InitialModel returned %T, want menuUnassigned", m)
	}

	if mu.list.Title != "Sneat.app - main menu" {
		t.Fatalf("list title = %q, want %q", mu.list.Title, "Sneat.app - main menu")
	}

	// Sanity check that view renders something
	view := mu.View()
	if view == "" {
		t.Fatalf("View() returned empty string")
	}
}

func TestMenuUnassigned_InitReturnsNil(t *testing.T) {
	m := InitialModel().(menuUnassigned)
	if cmd := m.Init(); cmd != nil {
		// execute to show behavior if non-nil
		_ = cmd()
		t.Fatalf("Init() = %v, want nil", cmd)
	}
}

func TestMenuUnassigned_Update_CtrlCQuits(t *testing.T) {
	m := InitialModel().(menuUnassigned)
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := model.(menuUnassigned); !ok {
		t.Fatalf("Update returned model %T, want menuUnassigned", model)
	}
	if cmd == nil {
		t.Fatalf("cmd is nil, want tea.Quit")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("cmd() returned %T, want tea.QuitMsg", msg)
	}
}

func TestMenuUnassigned_Update_WindowSizeSetsListSize(t *testing.T) {
	m := InitialModel().(menuUnassigned)
	// initial sizes
	w0, h0 := m.list.Width(), m.list.Height()

	hMargin, vMargin := docStyle.GetFrameSize()
	msg := tea.WindowSizeMsg{Width: 100, Height: 40}
	model, cmd := m.Update(msg)
	_ = cmd // may be nil; we're focused on size changes

	mu := model.(menuUnassigned)
	wantW := msg.Width - hMargin
	wantH := msg.Height - vMargin
	if mu.list.Width() != wantW || mu.list.Height() != wantH {
		t.Fatalf("list size = (%d,%d), want (%d,%d); initial was (%d,%d)", mu.list.Width(), mu.list.Height(), wantW, wantH, w0, h0)
	}
}
